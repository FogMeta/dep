package v1

import (
	"errors"
	"strconv"

	"github.com/FogMeta/libra-os/api"
	"github.com/FogMeta/libra-os/api/result"
	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/model/resp"
	"github.com/FogMeta/libra-os/module/log"
	"github.com/FogMeta/libra-os/service"
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
		api.ErrResponse(c, result.InternalError, err)
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
		api.ErrResponse(c, result.InternalError, err)
		return
	}
	api.Response(c, info)
}

func (api *DeploymentAPI) Deployments(c *gin.Context) {
	var req req.DeploymentQueryReq
	if err := api.ParseReq(c, &req, true); err != nil {
		return
	}
	uid := api.UID(c)
	if req.JobUUID == "" && req.SpaceUUID == "" {
		info, err := spaceService.DeploymentList(uid)
		if err != nil {
			api.ErrResponse(c, result.InternalError, err)
			return
		}
		api.Response(c, info)
		return
	}
	info, err := spaceService.Deployment(uid, req.JobUUID, req.SpaceUUID)
	if err != nil {
		api.ErrResponse(c, result.InternalError, err)
		return
	}
	api.Response(c, info)
}

func (api *DeploymentAPI) DeploymentList(c *gin.Context) {
	var req req.DeploymentQueryReq
	if err := api.ParseReq(c, &req, true); err != nil {
		return
	}
	page, size := req.PageNo, req.PageSize
	if size <= 0 {
		size = 10
	} else if size > 30 {
		size = 30
	}
	uid := api.UID(c)
	deployment := &model.Deployment{
		UID:    uid,
		Status: req.Status,
	}
	list, err := spaceService.Deployments(deployment, page*size, size)
	if err != nil {
		api.ErrResponse(c, result.InternalError, err)
		return
	}
	deployments := make([]*resp.DeploymentInfo, 0, len(list))
	for _, dp := range list {
		specs, _ := service.CfgMachine(dp.CfgName)
		deployments = append(deployments, &resp.DeploymentInfo{
			ID:             dp.ID,
			UID:            dp.UID,
			SpaceID:        dp.SpaceID,
			SpaceName:      dp.SpaceName,
			CfgName:        dp.CfgName,
			CfgSpecs:       specs,
			Duration:       dp.Duration,
			Region:         dp.Region,
			ResultURL:      dp.ResultURL,
			ProviderID:     dp.ProviderID,
			ProviderNodeID: dp.ProviderNodeID,
			Cost:           dp.Cost,
			Spent:          dp.Spent,
			ExpiredAt:      dp.ExpiredAt,
			EndedAt:        dp.EndedAt,
			Status:         dp.Status,
			StatusMsg:      dp.StatusMsg,
			Source:         dp.Source,
			CreatedAt:      dp.CreatedAt.Unix(),
		})
	}
	total, err := spaceService.Count(deployment)
	if err != nil {
		log.Error(err)
	}
	api.Response(c, resp.PageList{
		Total: int(total),
		List:  deployments,
	})
}

func (api *DeploymentAPI) DeploymentInfo(c *gin.Context) {
	uid := api.UID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		api.ErrResponse(c, result.InvalidPara, errors.New("invalid id"))
		return
	}
	info, err := spaceService.DeploymentInfo(uid, id)
	if err != nil {
		api.ErrResponse(c, result.InternalError, err)
		return
	}
	info.CfgSpecs, _ = service.CfgMachine(info.CfgName)
	api.Response(c, info)
}
