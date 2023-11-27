package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/model/resp"
	"github.com/FogMeta/libra-os/module/lagrange"
	"github.com/FogMeta/libra-os/module/log"
)

var errNotFoundKey = errors.New("not found api token")

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
	StatusStopped
)

var statusMap = map[string]int{
	"Running": StatusSuccess,
	"Stopped": StatusStopped,
}

type SpaceService struct {
	DBService
}

func (s *SpaceService) SpaceInfo(uid int, spaceURL string) (space *resp.SpaceResp, err error) {
	user, err := s.User(uid, true)
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
	user, err := s.User(uid, true)
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
	spaceUUID := req.SpaceUUID
	if result.Space.UUID != "" {
		spaceUUID = result.Space.UUID
	}
	dp := &model.Deployment{
		UID:       uid,
		SpaceID:   spaceUUID,
		SpaceName: req.SpaceName,
		CfgName:   req.CfgName,
		Paid:      req.Paid,
		Duration:  req.Duration,
		TxHash:    req.TxHash,
		ChainID:   req.ChainID,
		Region:    req.Region,
		StartIn:   req.StartIn,
		Source:    SourceLagrange,
		Status:    StatusReady,
	}
	if err = s.Insert(dp); err != nil {
		return
	}
	return &resp.DeploymentCreateResp{
		ID:                dp.ID,
		SpaceDeployResult: *result,
	}, nil
}

func (s *SpaceService) DeployStatus(uid int, spaceUUID string) (result *lagrange.DeployStatus, err error) {
	user, err := s.User(uid, true)
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
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	result, err = lagClient.WithAPIKey(user.APIKey).Deployment(jobUUID, spaceUUID)
	return
}

func (s *SpaceService) DeploymentList(uid int) (deployments []*lagrange.DeploymentAbstract, err error) {
	user, err := s.User(uid, true)
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
		Region:         dp.Region,
		ResultURL:      dp.ResultURL,
		ProviderID:     dp.ProviderID,
		ProviderNodeID: dp.ProviderNodeID,
		Cost:           dp.Cost,
		Status:         dp.Status,
		StatusMsg:      dp.StatusMsg,
		Source:         SourceLagrange,
		CreatedAt:      dp.CreatedAt.Unix(),
	}
	return
}

func (s *SpaceService) Deployments(deployment *model.Deployment, args ...int) (deployments []*model.Deployment, err error) {
	var wheres []string
	if deployment.Status == 1 {
		deployment.Status = 0
		wheres = append(wheres, fmt.Sprintf("status BETWEEN %d AND %d", StatusReady, StatusSuccess))
	}
	if err = s.Find(deployment, &deployments, wheres, args...); err != nil {
		return
	}

	for _, dp := range deployments {
		if dp.ResultURL != "" {
			continue
		}
		if err = s.LagrangeSync(dp); err != nil {
			log.Error(err)
		}
	}
	return
}

func (s *SpaceService) Count(deployment *model.Deployment) (count int64, err error) {
	var where string
	if deployment.Status == 1 {
		deployment.Status = 0
		where = fmt.Sprintf("status BETWEEN %d AND %d", StatusReady, StatusSuccess)
	}
	dm := s.DB().Model(deployment).Where(deployment)
	if where != "" {
		dm = dm.Where(where)
	}
	err = dm.Count(&count).Error
	return
}

func (s *DBService) LagrangeSync(dp *model.Deployment) (err error) {
	if dp.ResultURL != "" {
		return
	}
	user, err := s.User(dp.UID, true)
	if err != nil {
		return
	}
	expired := dp.CreatedAt.Add(time.Duration(dp.StartIn) * time.Second).Before(time.Now())
	if dp.JobID == "" {
		dp.JobID, err = lagClient.WithAPIKey(user.APIKey).JobID(dp.SpaceID)
		if err != nil {
			log.Error(err)
			if expired {
				return
			}
		}
		dp.StatusMsg = "Deploying"
	}
	result, err := lagClient.WithAPIKey(user.APIKey).Deployment(dp.JobID, dp.SpaceID)
	if err != nil {
		return
	}
	dp.StatusMsg = result.Status
	dp.ResultURL = result.ResultURL
	dp.ProviderNodeID = result.ProviderNodeID
	dp.Cost = result.ExpectedCost
	dp.StatusMsg = result.DeployStatus
	dp.Status = statusMap[result.DeployStatus]
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
