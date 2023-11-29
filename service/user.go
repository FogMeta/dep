package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/config"
	"github.com/FogMeta/libra-os/misc"
	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/model/resp"
	sdk "github.com/FogMeta/libra-os/module/lagrange"
	"github.com/FogMeta/libra-os/module/log"
	"github.com/FogMeta/libra-os/module/redis"
)

const (
	UserTypeWallet = iota + 1
	UserTypeEmail
)

type UserService struct {
	DBService
	JWTService
}

func (s *UserService) UserInfo(user *model.User) (res *resp.UserInfoResp, code int, err error) {
	if err := s.First(user); err != nil {
		return nil, result.UserNotFound, errors.New("user not found")
	}
	res = &resp.UserInfoResp{
		UID:    user.ID,
		Name:   user.Email,
		Email:  user.Email,
		Wallet: user.Wallet,
	}
	return
}

func (s *UserService) Register(req *req.UserCreateReq) (res *resp.UserResp, code int, err error) {
	if req.Type == UserTypeEmail {
		return s.RegisterWithEmail(&model.User{
			Email:    req.Email,
			Password: req.Password,
			Type:     req.Type,
		}, req.AuthCode)
	}
	return s.LoginWithWallet(&model.User{
		Type: req.Type,
	}, req.WalletToken)
}

// LoginWithWallet login or register without login
func (s *UserService) LoginWithWallet(user *model.User, token string) (res *resp.UserResp, code int, err error) {
	// get api key
	client := lagClient.WithToken(token)
	user.APIKey, err = client.APIToken()
	if err != nil {
		log.Error(err)
		return
	}

	// validate wallet token
	client = client.WithAPIKey(user.APIKey)
	user.Wallet, err = client.TokenWallet()
	if err != nil {
		log.Error(err)
		return
	}

	// query user
	dbu := &model.User{Wallet: user.Wallet}
	if err = s.First(dbu); err == nil {
		if dbu.APIKey != user.APIKey {
			// update api key
			dbu.APIKey = user.APIKey
			if err = s.Updates(dbu, "api_key"); err != nil {
				return
			}
		}
		return s.generateToken(dbu, false)
	}
	// create new user
	user.Status = 1
	if err = s.Insert(user); err != nil {
		return
	}
	return s.generateToken(user, false)
}

func (s *UserService) RegisterWithEmail(user *model.User, authCode string) (res *resp.UserResp, code int, err error) {
	// check auth code
	key := fmt.Sprintf(RedisEmailKeyFormat, user.Email)
	ctx := context.TODO()
	value, err := redis.RDB.Get(ctx, key).Result()
	if err != nil || value != authCode {
		return nil, result.UserEmailCodeInvalid, errors.New("invalid auth code")
	}
	redis.RDB.Del(ctx, key)

	// check repeat
	if err := s.First(&model.User{Email: user.Email}); err == nil {
		return nil, result.UserEmailRegistered, errors.New("current email is already registered")
	}
	// create new user
	user.Status = 1
	user.Password = misc.MD5(user.Password)
	if err = s.Insert(user); err != nil {
		return
	}
	return s.generateToken(user, true)
}

func (s *UserService) Login(req *req.UserLoginReq) (res *resp.UserResp, code int, err error) {
	if req.Type == UserTypeEmail {
		return s.LoginWithEmail(&model.User{
			Email:    req.Email,
			Password: req.Password,
			Type:     req.Type,
		})
	}
	return s.LoginWithWallet(&model.User{
		Type: req.Type,
	}, req.WalletToken)
}

func (s *UserService) LoginWithEmail(user *model.User) (res *resp.UserResp, code int, err error) {
	password := misc.MD5(user.Password)
	u := model.User{Email: user.Email}
	if err := s.First(&u); err != nil {
		return nil, result.UserEmailPasswordIncorrect, errors.New("email or password incorrect")
	}
	*user = u
	if user.Password != password {
		return nil, result.UserEmailPasswordIncorrect, errors.New("email or password incorrect")
	}
	return s.generateToken(user, false)
}

func (s *UserService) UpdatePassword(user *model.User, newPass string) (res *resp.UserResp, code int, err error) {
	password := misc.MD5(user.Password)
	u := model.User{ID: user.ID}
	if err := s.First(&u); err != nil {
		return nil, result.UserEmailPasswordIncorrect, errors.New("email or password incorrect")
	}
	*user = u
	if user.Password != password {
		return nil, result.UserEmailPasswordIncorrect, errors.New("email or password incorrect")
	}
	user.Password = newPass
	if err := s.Updates(user, "password"); err != nil {
		return nil, result.DBError, err
	}
	return s.generateToken(user, true)
}

func (s *UserService) ResetPassword(user *model.User, authCode string) (code int, err error) {
	// check auth code
	key := fmt.Sprintf(RedisEmailKeyFormat, user.Email)
	ctx := context.TODO()
	value, err := redis.RDB.Get(ctx, key).Result()
	if err != nil || value != authCode {
		return result.UserEmailCodeInvalid, errors.New("invalid auth code")
	}
	redis.RDB.Del(ctx, key)

	u := model.User{Email: user.Email}
	if err := s.First(&u); err != nil {
		return result.UserEmailPasswordIncorrect, errors.New("email or password incorrect")
	}
	u.Password = misc.MD5(user.Password)
	if err := s.Updates(u, "password"); err != nil {
		return result.DBError, err
	}
	*user = u
	return
}

func (s *UserService) generateToken(user *model.User, renew bool) (res *resp.UserResp, code int, err error) {
	token, err := s.GenerateToken(user.ID, true)
	if err != nil {
		return nil, result.RedisError, err
	}
	return &resp.UserResp{
		UID:    user.ID,
		Name:   user.Email,
		Email:  user.Email,
		Wallet: user.Wallet,
		Type:   user.Type,
		Token:  token,
	}, result.Success, nil
}

var lagClient *sdk.Client

func Init() {
	lagClient = &sdk.Client{
		Host:    config.Conf().Lagrange.Host,
		SDKHost: config.Conf().Lagrange.SDKHost,
	}
}
