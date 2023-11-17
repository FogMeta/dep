package v1

import (
	"github.com/FogMeta/libra-os/api"
	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/module/lagrange"
	"github.com/gin-gonic/gin"
)

type DeploymentAPI struct {
	api.BaseApi
}

func (api *DeploymentAPI) Deploy(c *gin.Context) {
	var req req.SpaceDeployReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	uid := api.UID(c)
	info, err := spaceService.Deploy(uid, &req)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}
	api.Response(c, info)
}

func (api *DeploymentAPI) DeployStatus(c *gin.Context) {
	var req req.DeploymentStatusReq
	if err := api.ParseReq(c, &req, true); err != nil {
		return
	}
	uid := api.UID(c)
	info, err := spaceService.DeployStatus(uid, req.SpaceUUID)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}
	api.Response(c, info)
}

func (api *DeploymentAPI) Deployments(c *gin.Context) {
	var req req.DeploymentQueryReq
	if err := api.ParseReq(c, &req); err != nil {
		return
	}
	uid := api.UID(c)
	if req.JobUUID == "" && req.SpaceUUID == "" {
		info, err := spaceService.DeploymentList(uid)
		if err != nil {
			api.ErrResponse(c, result.SpaceURLInvalid, err)
			return
		}
		if req.Status != "" {
			list := make([]*lagrange.DeploymentAbstract, 0, len(info))
			for _, deployment := range info {
				if deployment.Status == req.Status {
					list = append(list, deployment)
				}
			}
			info = list
		}
		api.Response(c, info)
		return
	}
	info, err := spaceService.Deployment(uid, req.JobUUID, req.SpaceUUID)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}
	api.Response(c, info)
}

func (api *DeploymentAPI) DeploymentList(c *gin.Context) {
	uid := api.UID(c)
	info, err := spaceService.DeploymentList(uid)
	if err != nil {
		api.ErrResponse(c, result.SpaceURLInvalid, err)
		return
	}
	api.Response(c, info)
}
