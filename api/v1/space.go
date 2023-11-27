package v1

import (
	"errors"

	"github.com/FogMeta/libra-os/api"
	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/service"
	"github.com/gin-gonic/gin"
)

var spaceService = new(service.SpaceService)

type SpaceAPI struct {
	api.BaseApi
}

func (api *SpaceAPI) SpaceInfo(c *gin.Context) {
	var req req.SpaceReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	uid := api.UID(c)
	user, err := spaceService.User(uid)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}
	info, err := spaceService.SpaceInfo(uid, req.URL)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}
	if user.Wallet != info.Wallet {
		api.ErrResponse(c, result.SpaceWalletNotMatch, errors.New("space must be your own, please fork to your space and retry"))
		return
	}
	api.Response(c, info)
}
