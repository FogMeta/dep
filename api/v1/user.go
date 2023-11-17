package v1

import (
	"github.com/FogMeta/libra-os/api"
	"github.com/FogMeta/libra-os/misc"
	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/service"
	"github.com/gin-gonic/gin"
)

var userService = new(service.UserService)

type UserApi struct {
	api.BaseApi
}

func (api *UserApi) ResetPassword(c *gin.Context) {
	var req req.UserResetPasswordReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}
	code, err := userService.ResetPassword(user, req.AuthCode)
	if err != nil {
		api.ErrResponse(c, code, err)
		return
	}
	api.Response(c, nil)
}

func (api *UserApi) Register(c *gin.Context) {
	var req req.UserCreateReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	resp, code, err := userService.Register(&req)
	if err != nil {
		api.ErrResponse(c, code, err)
		return
	}
	api.Response(c, resp)
}

func (api *UserApi) Login(c *gin.Context) {
	var req req.UserLoginReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	resp, code, err := userService.Login(&req)
	if err != nil {
		api.ErrResponse(c, code, err)
		return
	}
	api.Response(c, resp)
}

func (api *UserApi) UpdatePassword(c *gin.Context) {
	var req req.UserUpdatePasswordReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	user := &model.User{
		ID:       api.UID(c),
		Password: req.OldPassword,
	}
	resp, code, err := userService.UpdatePassword(user, misc.MD5(req.Password))
	if err != nil {
		api.ErrResponse(c, code, err)
		return
	}
	api.Response(c, resp)
}

func (api *UserApi) UserInfo(c *gin.Context) {
	user := &model.User{
		ID: api.UID(c),
	}
	resp, code, err := userService.UserInfo(user)
	if err != nil {
		api.ErrResponse(c, code, err)
		return
	}
	api.Response(c, resp)
}
