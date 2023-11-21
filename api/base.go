package api

import (
	"net/http"
	"net/http/httputil"

	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/module/log"
	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

func (a *BaseApi) UID(c *gin.Context) int {
	return c.GetInt("uid")
}

func (api *BaseApi) ParseReq(c *gin.Context, receiverPointer any, query ...bool) error {
	body, _ := httputil.DumpRequest(c.Request, true)
	log.Info(string(body))
	if len(query) > 0 && query[0] {
		if err := c.BindQuery(receiverPointer); err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
	if err := c.BindJSON(receiverPointer); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (api *BaseApi) Response(c *gin.Context, data any) {
	api.response(c, result.Result{
		Code: result.Success,
		Data: data,
		Msg:  "success",
	})
}

func (api *BaseApi) response(c *gin.Context, data any, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	c.JSON(code, data)
}

func (api *BaseApi) ErrResponse(c *gin.Context, code int, err error, statusCode ...int) {
	status := http.StatusBadRequest
	if len(statusCode) > 0 {
		status = statusCode[0]
	}
	res := result.Result{
		Code: code,
		Msg:  err.Error(),
	}
	api.response(c, res, status)
}
