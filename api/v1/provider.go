package v1

import (
	"errors"
	"net/http"

	"github.com/FogMeta/libra-os/api"
	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/service"
	"github.com/gin-gonic/gin"
)

var providerService = new(service.ProviderService)

type ProviderAPI struct {
	api.BaseApi
}

func (api *ProviderAPI) ProviderList(c *gin.Context) {
	var req req.ProviderQuery
	if err := api.ParseReq(c, &req, true); err != nil {
		return
	}
	uid := api.UID(c)
	info, err := providerService.ProviderList(uid, req.Region)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}

	api.Response(c, info)
}

func (api *ProviderAPI) Provider(c *gin.Context) {
	uuid, _ := c.Params.Get("uuid")
	if uuid == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid uuid"))
	}
	uid := api.UID(c)

	info, err := providerService.Provider(uid, uuid)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}

	api.Response(c, info)
}

func (api *ProviderAPI) ProviderDistribution(c *gin.Context) {
	var req req.ProviderQuery
	if err := api.ParseReq(c, &req, true); err != nil {
		return
	}
	uid := api.UID(c)
	info, err := providerService.ProviderDistribution(uid, req.Region)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}

	api.Response(c, info)
}

func (api *ProviderAPI) Resources(c *gin.Context) {
	uid := api.UID(c)
	info, err := providerService.ResourceSummary(uid)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}

	api.Response(c, info)
}
