package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FogMeta/libra-os/misc"
	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/model/req"
	"github.com/FogMeta/libra-os/model/resp"
	"github.com/FogMeta/libra-os/module/lagrange"
	"github.com/FogMeta/libra-os/module/log"
)

var (
	errNotFoundKey       = errors.New("not found api token")
	errSpaceDeploying    = errors.New("Space does not have a corresponding job")
	errSpaceDeployFailed = errors.New("Space does not have a corresponding Task")
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
	StatusStopped
)

var statusMap = map[string]int{
	"Running":                 StatusSuccess,
	"Stopped":                 StatusStopped,
	"Pending":                 StatusTransaction,
	"Waiting for transaction": StatusTransaction,
	"Assigning to provider":   StatusAssignProvider,
	"Complete":                StatusStopped,
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
	status := StatusReady
	statusMsg := result.Space.Status
	if result.Payment != nil && result.Payment.Status == "Pending" {
		status = StatusTransaction
		statusMsg = "Waiting for transaction"
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
		Status:    status,
		StatusMsg: statusMsg,
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
		log.Error(err)
	}
	return
}

func (s *SpaceService) Deployment(uid int, jobUUID, spaceUUID string) (result *lagrange.Deployment, err error) {
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	result, err = lagClient.WithAPIKey(user.APIKey).Deployment(jobUUID, spaceUUID)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *SpaceService) DeploymentList(uid int) (deployments []*lagrange.DeploymentAbstract, err error) {
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	deployments, err = lagClient.WithAPIKey(user.APIKey).DeploymentList()
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *SpaceService) DeploymentInfo(uid int, id int) (deployment *resp.DeploymentInfo, err error) {
	dp := &model.Deployment{UID: uid, ID: id}
	if err = s.First(dp); err != nil {
		return
	}
	if err = s.LagrangeSync(dp); err != nil {
		log.Error(err)
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
		LastResultURL:  dp.LastResultURL,
		ProviderID:     dp.ProviderID,
		ProviderNodeID: dp.ProviderNodeID,
		Cost:           dp.Cost,
		Spent:          dp.Spent,
		ExpiredAt:      dp.ExpiredAt,
		EndedAt:        dp.EndedAt,
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
		log.Error(err)
		return
	}

	for _, dp := range deployments {
		if err := s.LagrangeSync(dp); err != nil {
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
	if dp.Status == StatusStopped {
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
			if err == errSpaceDeploying {
				dp.StatusMsg = "Deploying"
				return nil
			}
			if err == errSpaceDeployFailed {
				dp.Status = StatusFailed
				dp.StatusMsg = err.Error()
				return nil
			}
			if expired {
				log.Infof("deployment %d space %s deploy expired", dp.ID, dp.SpaceID)
				dp.Status = StatusFailed
				dp.StatusMsg = err.Error()
				return nil
			}
			return err
		}
		dp.StatusMsg = "Deploying"
		s.Updates(dp, "job_id")
	}
	deployment := *dp
	result, err := lagClient.WithAPIKey(user.APIKey).Deployment(dp.JobID, dp.SpaceID)
	if err != nil {
		log.Error(err)
		return
	}
	if result.ResultURL != "" && dp.LastResultURL == "" {
		dp.LastResultURL, err = lagrange.ResultURL(result.ResultURL)
		if err != nil {
			log.Error(err)
		}
	}
	dp.StatusMsg = result.Status
	dp.ResultURL = result.ResultURL
	dp.ProviderNodeID = result.ProviderNodeID
	dp.Region = result.Region
	dp.Cost = result.ExpectedCost
	dp.Spent = int(result.Spent)
	dp.StatusMsg = result.DeployStatus
	dp.Status = statusMap[result.DeployStatus]
	dp.ExpiredAt = result.ExpiresAt
	if result.EndedAt != nil {
		endedAt, err := strconv.ParseInt(*result.EndedAt, 10, 64)
		if err != nil {
			log.Error(err)
		} else {
			dp.EndedAt = endedAt
		}
	}

	// update db
	values, err := misc.CompareStructValues(deployment, dp, "gorm", "id", "created_at", "updated_at")
	if err != nil {
		log.Error(err)
		return
	}
	if len(values) > 0 {
		log.Info("update values:", values)
		return s.DB().Model(dp).Updates(values).Error
	}
	return nil
}
