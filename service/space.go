package service

import (
	"strings"

	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/model/resp"
	"github.com/FogMeta/libra-os/module/lagrange"
)

type SpaceService struct {
	DBService
}

func (s *SpaceService) SpaceInfo(uid int, spaceURL string) (space *resp.SpaceResp, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	space = new(resp.SpaceResp)
	space.UUID, err = lagClient.WithAPIKey(user.APIKey).SpaceUUID(spaceURL)
	if err != nil {
		return
	}
	list := strings.Split(strings.TrimPrefix(spaceURL, "https://"), "/")
	if len(list) >= 5 {
		space.Wallet = list[2]
		space.Name = list[3]
	}
	return
}

func (s *SpaceService) Deploy(uid int, req *req.SpaceDeployReq) (result *lagrange.SpaceDeployResult, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	return lagClient.WithAPIKey(user.APIKey).Deploy(&lagrange.SpaceDeployReq{
		SpaceName: req.SpaceName,
		Paid:      req.Paid,
		Duration:  req.Duration,
		TxHash:    req.TxHash,
		ChainID:   req.ChainID,
		CfgName:   req.CfgName,
		Region:    req.Region,
		StartIn:   req.StartIn,
	})
}

func (s *SpaceService) DeployStatus(uid int, spaceUUID string) (result *lagrange.DeployStatus, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	result, err = lagClient.WithAPIKey(user.APIKey).DeployStatus(spaceUUID)
	if err != nil {
		return
	}
	return
}

func (s *SpaceService) Deployment(uid int, jobUUID, spaceUUID string) (result *lagrange.Deployment, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	result, err = lagClient.WithAPIKey(user.APIKey).Deployment(jobUUID, spaceUUID)
	if err != nil {
		return
	}
	return
}

func (s *SpaceService) DeploymentList(uid int) (deployments []*lagrange.DeploymentAbstract, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	deployments, err = lagClient.WithAPIKey(user.APIKey).DeploymentList()
	if err != nil {
		return
	}
	return
}
