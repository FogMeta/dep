package resp

type DeploymentCreateResp struct {
	JobUUID string `json:"job_uuid"`
}

type DeploymentListResp struct {
	SpaceName string `json:"space_name"`
	SpaceUUID string `json:"space_uuid"`
	Status    int    `json:"status"`
}

type DeploymentResp struct {
	ProviderNodeID string `json:"provider_node_id"`
	Region         string `json:"region"`
	DeployStatus   int    `json:"deploy_status"`
	ResultURL      string `json:"result_url"`
	EndedAt        int    `json:"ended_at"`
	Status         string `json:"status"`
	ExpiredAt      int    `json:"expired_at"`
	Balance        string `json:"balance"`
	Cost           string `json:"cost"`
}
