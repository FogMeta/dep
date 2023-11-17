package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FogMeta/libra-os/misc"
	"github.com/FogMeta/libra-os/module/log"
	"github.com/FogMeta/libra-os/module/redis"
	"github.com/golang-jwt/jwt/v5"
)

const (
	redisUserSecretKeyFormat = "libra:jwt:uid:%d:secret"
	redisUserTokenFresh      = "libra:jwt:uid:%d:token"
)

type JWTService struct{}

func (s *JWTService) Validate(token string) (uid int, newToken string, err error) {
	claims := new(Claims)
	if err = claims.ParseToken(token); err != nil {
		log.Error(err)
		return
	}
	uid = claims.User.ID
	if time.Until(claims.ExpiresAt.Time) < tokenExpireDuration/2 {
		ctx := context.Background()
		key := fmt.Sprintf(redisUserTokenFresh, uid)
		redis.RDB.SetNX(ctx, key, 1, time.Second*10)
		defer redis.RDB.Del(ctx, key)
		newToken, err = claims.User.GenerateToken(false)
		return
	}
	return
}

func (s *JWTService) GenerateToken(id int, renew bool) (token string, err error) {
	user := &UserInfo{ID: id}
	return user.GenerateToken(renew)
}

type Claims struct {
	jwt.RegisteredClaims
	User *UserInfo `json:"user"`
}

func (claims *Claims) ParseToken(token string) error {
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return claims.User.Secret()
	})
	return err
}

func (claims *Claims) RenewToken() (string, error) {
	if time.Since(claims.ExpiresAt.Time) < 10*time.Minute {
		return claims.User.GenerateToken(false)
	}
	return "", errors.New("token is expired")
}

const (
	tokenExpireDuration = time.Hour * 2

	secretLength         = 16
	secretExpireDuration = time.Hour * 24 * 30
)

type UserInfo struct {
	ID int `json:"id"`
}

func (user *UserInfo) GenerateToken(renew bool) (string, error) {
	var secret []byte
	if renew {
		secret = []byte(misc.RandomString(secretLength, misc.CharsetAll))
		if err := redis.RDB.Set(context.TODO(), fmt.Sprintf(redisUserSecretKeyFormat, user.ID), secret, secretExpireDuration).Err(); err != nil {
			return "", err
		}
	} else {
		var err error
		secret, err = user.Secret()
		if err != nil {
			return "", err
		}
	}
	expirationTime := time.Now().Add(tokenExpireDuration)
	claims := &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (user *UserInfo) Secret() ([]byte, error) {
	if user == nil {
		return nil, errors.New("invalid nil user")
	}
	key, err := redis.RDB.Get(context.TODO(), fmt.Sprintf(redisUserSecretKeyFormat, user.ID)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(key), nil
}
