package lagrange

import (
	"net/url"

	"github.com/FogMeta/libra-os/misc"
	"github.com/FogMeta/libra-os/module/log"
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
	result.SpaceUUID = &uuid
	err = client.postForm(methodSpaceUUID, data, &result)
	return
}

type spaceUUIDResult struct {
	SpaceUUID *string `json:"space_uuid"`
}

func (client *Client) Deploy(req *SpaceDeployReq) (status *SpaceDeployResult, err error) {
	values, err := misc.EncodeStructValues(req, "json", "&")
	if err != nil {
		return
	}
	log.Info("values: ", values)
	var result SpaceDeployResult
	if err = client.postForm(methodSpaceDeploy, values, &result); err != nil {
		return
	}
	return &result, nil
}

type SpaceDeployReq struct {
	SpaceName string `json:"space_name"`
	Paid      string `json:"paid"`
	Duration  int    `json:"duration"`
	TxHash    string `json:"tx_hash"`
	ChainID   string `json:"chain_id"`
	CfgName   string `json:"cfg_name"`
	Region    string `json:"region"`
	StartIn   int    `json:"start_in"`
}

type SpaceDeployResult struct {
	Space             Space              `json:"space"`
	Task              Task               `json:"task"`
	DeploymentRequest *DeploymentRequest `json:"deployment_request"`
	Payment           *Payment           `json:"payment"`
}

type DeploymentRequest struct {
	ChainID         int64        `json:"chain_id"`
	CreatedAt       string       `json:"created_at"`
	Order           *DeployOrder `json:"order"`
	SpaceName       string       `json:"space_name"`
	TransactionHash string       `json:"transaction_hash"`
	UpdatedAt       string       `json:"updated_at"`
}

type DeployOrder struct {
	Config    *DeplyConfig `json:"config"`
	CreatedAt string       `json:"created_at"`
	Duration  int64        `json:"duration"`
	EndedAt   interface{}  `json:"ended_at"`
	OrderType interface{}  `json:"order_type"`
	SpaceName string       `json:"space_name"`
	StartedAt interface{}  `json:"started_at"`
	UpdatedAt string       `json:"updated_at"`
}

type DeplyConfig struct {
	Description  string      `json:"description"`
	Hardware     interface{} `json:"hardware"`
	HardwareType string      `json:"hardware_type"`
	Memory       int64       `json:"memory"`
	Name         string      `json:"name"`
	PricePerHour float64     `json:"price_per_hour"`
	Vcpu         int64       `json:"vcpu"`
}

type Payment struct {
	Amount           string      `json:"amount"`
	ChainID          int64       `json:"chain_id"`
	CreatedAt        string      `json:"created_at"`
	DeniedReason     interface{} `json:"denied_reason"`
	Order            DeployOrder `json:"order"`
	RefundReason     interface{} `json:"refund_reason"`
	RefundableAmount interface{} `json:"refundable_amount"`
	Status           string      `json:"status"`
	Token            interface{} `json:"token"`
	TransactionHash  string      `json:"transaction_hash"`
	UpdatedAt        string      `json:"updated_at"`
}

type Space struct {
	ActiveOrder    ActiveOrder `json:"activeOrder"`
	CreatedAt      string      `json:"created_at"`
	ExpirationTime string      `json:"expiration_time"`
	IsPublic       int64       `json:"is_public"`
	LastStopReason interface{} `json:"last_stop_reason"`
	License        string      `json:"license"`
	Likes          int64       `json:"likes"`
	Name           string      `json:"name"`
	SDK            string      `json:"sdk"`
	Status         string      `json:"status"`
	TaskUUID       string      `json:"task_uuid"`
	UpdatedAt      string      `json:"updated_at"`
	UUID           string      `json:"uuid"`
}

type ActiveOrder struct {
	Config    Config      `json:"config"`
	CreatedAt string      `json:"created_at"`
	Duration  int64       `json:"duration"`
	EndedAt   interface{} `json:"ended_at"`
	OrderType interface{} `json:"order_type"`
	SpaceName string      `json:"space_name"`
	StartedAt int64       `json:"started_at"`
	UpdatedAt string      `json:"updated_at"`
}

type Config struct {
	Description  string      `json:"description"`
	Hardware     interface{} `json:"hardware"`
	HardwareType string      `json:"hardware_type"`
	Memory       int64       `json:"memory"`
	Name         string      `json:"name"`
	PricePerHour float64     `json:"price_per_hour"`
	Vcpu         int64       `json:"vcpu"`
}

type Task struct {
	CreatedAt     string      `json:"created_at"`
	LeadingJobID  interface{} `json:"leading_job_id"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	TaskDetailCid string      `json:"task_detail_cid"`
	UpdatedAt     string      `json:"updated_at"`
	UUID          string      `json:"uuid"`
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

func (client *Client) Deployment(jobUUID, spaceUUID string) (deployment *Deployment, err error) {
	var result Deployment
	if err = client.get(methodSpaceDeployment+"/"+jobUUID+"/"+spaceUUID, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type Deployment struct {
	DeployStatus   string  `json:"deploy_status"`
	EndedAt        *string `json:"ended_at"`
	ExpectedCost   string  `json:"expected_cost"`
	ExpiresAt      int64   `json:"expires_at"`
	ProviderNodeID string  `json:"provider_node_id"`
	Region         string  `json:"region"`
	ResultURL      string  `json:"result_url"`
	Spent          int64   `json:"spent"`
	Status         string  `json:"status"`
}

func (client *Client) DeploymentList() (deployments []*DeploymentAbstract, err error) {
	err = client.get(methodDeploymentList, nil, &deployments)
	return
}

type DeploymentAbstract struct {
	JobUUID   string `json:"job_uuid"`
	SpaceName string `json:"space_name"`
	SpaceUUID string `json:"space_uuid"`
	Status    string `json:"status"`
}
