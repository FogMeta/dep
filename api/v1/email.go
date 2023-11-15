package v1

import (
	"github.com/FogMeta/libra-os/api"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/service"
	"github.com/gin-gonic/gin"
)

var emailService = new(service.EmailService)

type EmailApi struct {
	api.BaseApi
}

func (api *EmailApi) Send(c *gin.Context) {
	var req req.EmailReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	code, err := emailService.SendEmail(req.Email)
	if err != nil {
		api.ErrResponse(c, code, err)
		return
	}
	api.Response(c, nil)
}
