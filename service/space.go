package service

import (
	"strings"
	"time"

	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/model/resp"
	"github.com/FogMeta/libra-os/module/lagrange"
	"github.com/FogMeta/libra-os/module/log"
)

const (
	SourceLagrange = iota + 1
)

const (
	StatusFailed = iota + 10
	StatusReady
	StatusTransaction
	StatusAssignProvider
	StatusDeploying
	StatusSuccess
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

func (s *SpaceService) Deploy(uid int, req *req.SpaceDeployReq) (deployment *resp.DeploymentCreateResp, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	result, err := lagClient.WithAPIKey(user.APIKey).Deploy(&lagrange.SpaceDeployReq{
		SpaceName: req.SpaceName,
		Paid:      req.Paid,
		Duration:  req.Duration,
		TxHash:    req.TxHash,
		ChainID:   req.ChainID,
		CfgName:   req.CfgName,
		Region:    req.Region,
		StartIn:   req.StartIn,
	})
	if err != nil {
		return
	}
	dp := &model.Deployment{
		UID:       uid,
		SpaceID:   result.Space.UUID,
		SpaceName: req.SpaceName,
		CfgName:   req.CfgName,
		Paid:      req.Paid,
		Duration:  req.Duration,
		TxHash:    req.TxHash,
		ChainID:   req.ChainID,
		Region:    req.Region,
		StartIn:   req.StartIn,
		Source:    SourceLagrange,
	}
	if err = s.Insert(dp); err != nil {
		return
	}
	deployment.ID = dp.ID
	return
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
	return
}

func (s *SpaceService) DeploymentList(uid int) (deployments []*lagrange.DeploymentAbstract, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	deployments, err = lagClient.WithAPIKey(user.APIKey).DeploymentList()
	return
}

func (s *SpaceService) DeploymentInfo(uid int, id int) (deployment *resp.DeploymentInfo, err error) {
	dp := &model.Deployment{UID: uid, ID: id}
	if err = s.First(dp); err != nil {
		return
	}
	if err = s.LagrangeSync(dp); err != nil {
		return
	}
	deployment = &resp.DeploymentInfo{
		ID:             id,
		UID:            uid,
		SpaceID:        dp.SpaceID,
		SpaceName:      dp.SpaceName,
		CfgName:        dp.CfgName,
		Duration:       dp.Duration,
		Region:         dp.ChainID,
		ResultURL:      dp.ResultURL,
		ProviderID:     dp.ProviderID,
		ProviderNodeID: dp.ProviderNodeID,
		Cost:           dp.Cost,
		Status:         dp.Status,
		StatusMsg:      dp.StatusMsg,
		Source:         SourceLagrange,
	}
	return
}

func (s *DBService) LagrangeSync(dp *model.Deployment) (err error) {
	if dp.ResultURL != "" {
		return
	}
	user, err := s.User(dp.UID)
	if err != nil {
		return
	}
	expired := dp.CreatedAt.Add(time.Duration(dp.Duration) * time.Second).Before(time.Now())
	if dp.JobID == "" {
		dp.JobID, err = lagClient.WithAPIKey(user.APIKey).JobID(dp.SpaceID)
		if err != nil {
			log.Error(err)
			if expired {
				return
			}
		}
		dp.StatusMsg = "deployment is deploying"
		return
	}
	result, err := lagClient.WithAPIKey(user.APIKey).Deployment(dp.JobID, dp.SpaceID)
	if err != nil {
		return
	}
	dp.StatusMsg = result.Status
	dp.ResultURL = result.ResultURL
	dp.ProviderNodeID = result.ProviderNodeID
	dp.Cost = result.ExpectedCost
	if result.ResultURL != "" {
		dp.Status = StatusSuccess
	}
	// update db
	s.DB().Model(dp).Updates(&model.Deployment{
		StatusMsg:      result.Status,
		ResultURL:      result.ResultURL,
		ProviderNodeID: result.ProviderNodeID,
		Cost:           result.ExpectedCost,
		Status:         dp.Status,
	})
	return
}
