package resp

import "github.com/FogMeta/libra-os/module/lagrange"

type DeploymentCreateResp struct {
	ID int `json:"id"`
	lagrange.SpaceDeployResult
}

type DeploymentAbstract struct {
	ID int `json:"id"`
	lagrange.DeploymentAbstract
}

type Deployment struct {
	ID int `json:"id"`
	lagrange.Deployment
}

type DeploymentInfo struct {
	ID             int    `json:"id"`
	UID            int    `json:"uid"`
	SpaceID        string `json:"space_id"`
	SpaceName      string `json:"space_name"`
	CfgName        string `json:"cfg_name"`
	CfgSpecs       any    `json:"cfg_specs"`
	Duration       int    `json:"duration"`
	Region         string `json:"region"`
	ResultURL      string `json:"result_url"`
	ProviderID     string `json:"provider_id"`
	ProviderNodeID string `json:"provider_node_id"`
	Cost           string `json:"cost"`
	Spent          int    `json:"spent"`
	ExpiredAt      int64  `json:"expired_at"`
	EndedAt        int64  `json:"ended_at"`
	Status         int    `json:"status"`
	StatusMsg      string `json:"status_msg"`
	Source         int    `json:"source"`
	CreatedAt      int64  `json:"created_at"`
}

type PageList struct {
	Total int `json:"total"`
	List  any `json:"list"`
}
