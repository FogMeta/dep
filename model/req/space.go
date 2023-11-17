package req

type SpaceReq struct {
	URL string `json:"url" form:"url"`
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

type DeploymentStatusReq struct {
	SpaceUUID string `json:"space_uuid"  form:"space_uuid"`
}

type DeploymentQueryReq struct {
	JobUUID   string `json:"job_uuid"`
	SpaceUUID string `json:"space_uuid"`
	Status    string `json:"status"`
}

type DeploymentConfig struct {
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	SpaceName string `json:"space_name"`
	Region    string `json:"region"`
	StartIn   int    `json:"start_in"`
}
