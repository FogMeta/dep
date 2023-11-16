package lagrange

import (
	"net/url"

	"github.com/FogMeta/libra-os/misc"
)

const (
	methodSpaceUUID       = "/space_uuid"
	methodSpaceDeploy     = "/deploy_space"
	methodJobUUID         = "/job_uuid"
	methodSpaceDeployment = "/space_deployment"
	methodDeploymentList  = "/user_deployment_list"
)

func (client *Client) SpaceUUID(spaceURL string) (uuid string, err error) {
	data := url.Values{}
	data.Set("space_url", spaceURL)
	var result spaceUUIDResult
	if err = client.postForm(methodSpaceUUID, data, &result); err != nil {
		return
	}
	return result.UUID, nil
}

type spaceUUIDResult struct {
	UUID string `json:"uuid"`
}

func (client *Client) Deploy(req *SpaceDeployReq) (status *SpaceDeployResult, err error) {
	values, err := misc.EncodeStructValues(req, "json")
	if err != nil {
		return
	}
	var result SpaceDeployResult
	if err = client.postForm(methodSpaceDeploy, values, &result); err != nil {
		return
	}
	return &result, nil
}

type SpaceDeployReq struct {
	SpaceName string `json:"space_name"`
	Paid      int    `json:"paid"`
	Duration  int    `json:"duration"`
	TxHash    string `json:"tx_hash"`
	ChainID   string `json:"chain_id"`
	CfgName   string `json:"cfg_name"`
	Region    string `json:"region"`
	StartIn   int    `json:"start_in"`
}

type SpaceDeployResult struct {
}

func (client *Client) DeployStatus(uuid string) (status *DeployStatus, err error) {
	var result DeployStatus
	if err = client.get(methodSpaceDeploy+"/"+uuid, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type DeployStatus struct {
	DeployStatus   string `json:"deploy_status"`
	JobUUID        string `json:"job_uuid"`
	ProviderNodeID string `json:"provider_node_id"`
	ResultURL      string `json:"result_url"`
}

func (client *Client) JobID(spaceUUID string) (uuid string, err error) {
	var result JobIDResult
	result.JobUUID = &uuid
	err = client.get(methodJobUUID+"/"+spaceUUID, nil, &result)
	return
}

type JobIDResult struct {
	JobUUID *string `json:"job_uuid"`
}

func (client *Client) Deployment(obUUID, spaceUUID string) (deployment *Deployment, err error) {
	var result Deployment
	if err = client.get(methodJobUUID+"/"+spaceUUID, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type Deployment struct {
	DeployStatus   string      `json:"deploy_status"`
	EndedAt        interface{} `json:"ended_at"`
	ExpectedCost   string      `json:"expected_cost"`
	ExpiresAt      int64       `json:"expires_at"`
	ProviderNodeID string      `json:"provider_node_id"`
	Region         string      `json:"region"`
	ResultURL      string      `json:"result_url"`
	Spent          int64       `json:"spent"`
	Status         string      `json:"status"`
}

func (client *Client) DeploymentList(jobUUID, spaceUUID string) (deployments []*DeploymentAbstract, err error) {
	err = client.get(methodSpaceDeployment+"/"+jobUUID+"/"+spaceUUID, nil, &deployments)
	return
}

type DeploymentAbstract struct {
	JobUUID   string `json:"job_uuid"`
	SpaceName string `json:"space_name"`
	SpaceUUID string `json:"space_uuid"`
	Status    string `json:"status"`
}
