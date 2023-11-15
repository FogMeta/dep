package service

import (
	"context"
	"fmt"
	"time"

	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/config"
	"github.com/FogMeta/libra-os/misc"
	"github.com/FogMeta/libra-os/module/redis"
	"gopkg.in/gomail.v2"
)

const RedisEmailKeyFormat = "ark:email:%s"

type EmailService struct{}

func (s *EmailService) SendEmail(email string) (code int, err error) {
	authCode := misc.RandomString(6, misc.CharsetNum)
	if err = redis.RDB.Set(context.TODO(), fmt.Sprintf(RedisEmailKeyFormat, email), authCode, 10*time.Minute).Err(); err != nil {
		return result.RedisError, err
	}
	conf := config.Conf().Email
	emailer := &Emailer{
		Sender:   conf.Email,
		Password: conf.Password,
		Host:     conf.Host,
		Port:     conf.Port,
	}
	err = emailer.Send(conf.Subject, fmt.Sprintf(conf.ContentTemplate, code), email)
	if err != nil {
		return result.UserEmailSendFailed, err
	}
	return
}

type Emailer struct {
	Sender   string
	Password string
	Host     string
	Port     int
}

func (e *Emailer) Send(subject, msg string, receivers ...string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Sender)
	m.SetHeader("To", receivers...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", msg)
	d := gomail.NewDialer(e.Host, e.Port, e.Sender, e.Password)
	return d.DialAndSend(m)
}
