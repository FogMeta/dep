package lagrange

import "net/url"

const (
	methodProviderList         = "/cp_list"
	methodProviderDistribution = "/cp_distribution"
	methodProviderDetail       = "/cp_detail"
	methodProviderAvailable    = "/cp_available"
	methodProviderMachines     = "/cp/machines"
	methodProviderDashboard    = "/cp/dashboard"
)

func (client *Client) ProviderList(region string) (providers []*Provider, err error) {
	var values url.Values
	if region != "" {
		values = url.Values{}
		values.Set("region", region)
	}
	if err = client.postForm(methodProviderList, values, &providers); err != nil {
		return
	}
	return
}

func (client *Client) Provider(providerID string) (provider *Provider, err error) {
	var result Provider
	if err = client.get(methodProviderDetail+"/"+providerID, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type Provider struct {
	ActiveDeployments int64         `json:"active_deployments"`
	Address           string        `json:"address"`
	Free              Resource      `json:"free"`
	Level             int64         `json:"level"`
	Location          string        `json:"location"`
	Name              string        `json:"name"`
	NodeID            string        `json:"node_id"`
	Specs             ResourceSpecs `json:"specs"`
	Total             Resource      `json:"total"`
	Uptime            int           `json:"uptime"`
	Used              Resource      `json:"used"`
}

type ResourceSpecs struct {
	CPU             string        `json:"cpu"`
	CPUArchitecture string        `json:"cpu_architecture"`
	Gpus            []interface{} `json:"gpus"`
}

func (client *Client) ProviderDistribution(region string) (distributions []*Distribution, err error) {
	err = client.get(methodProviderDistribution, nil, &distributions)
	return
}

type Distribution struct {
	City  string    `json:"city"`
	Value []float64 `json:"value"`
}

func (client *Client) ResourceSummary() (resource *ProviderResource, err error) {
	var result ProviderResource
	if err = client.get(methodProviderAvailable, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type ProviderResource struct {
	ActiveProviderCount int64    `json:"activeProviderCount"`
	Free                Resource `json:"free"`
	Total               Resource `json:"total"`
	Used                Resource `json:"used"`
}

type Resource struct {
	CPU     int `json:"cpu"`
	GPU     int `json:"gpu"`
	Memory  int `json:"memory"`
	Storage int `json:"storage"`
}

func (client *Client) Machines() (data *HardwareData, err error) {
	var result HardwareData
	if err = client.get(methodProviderMachines, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type HardwareData struct {
	Hardware []*Hardware `json:"hardware"`
}

type Hardware struct {
	HardwareDescription string   `json:"hardware_description"`
	HardwareID          int64    `json:"hardware_id"`
	HardwareName        string   `json:"hardware_name"`
	HardwarePrice       string   `json:"hardware_price"`
	HardwareStatus      string   `json:"hardware_status"`
	HardwareType        string   `json:"hardware_type"`
	Region              []string `json:"region"`
}

func (client *Client) Dashboard() (dashboard *Dashboard, err error) {
	var result Dashboard
	if err = client.get(methodProviderDashboard, nil, &result); err != nil {
		return
	}
	return &result, nil
}

type Dashboard struct {
	ActiveApplications int64           `json:"active_applications"`
	MapInfo            []*MapInfo      `json:"map_info"`
	Providers          []*ProviderInfo `json:"providers"`
	TotalCPU           int64           `json:"total_cpu"`
	TotalDeployments   int64           `json:"total_deployments"`
	TotalGPU           int64           `json:"total_gpu"`
	TotalMemory        int64           `json:"total_memory"`
	TotalProviders     int64           `json:"total_providers"`
	TotalStorage       int64           `json:"total_storage"`
	TotalUsedCPU       int64           `json:"total_used_cpu"`
	TotalUsedGPU       int64           `json:"total_used_gpu"`
	TotalUsedMemory    int64           `json:"total_used_memory"`
	TotalUsedStorage   int64           `json:"total_used_storage"`
	TotalUsedVcpu      int64           `json:"total_used_vcpu"`
	TotalVcpu          int64           `json:"total_vcpu"`
}

type MapInfo struct {
	City  string    `json:"city"`
	Value []float64 `json:"value"`
}

type ProviderInfo struct {
	City             string           `json:"city"`
	ComputerProvider ComputerProvider `json:"computer_provider"`
	Country          interface{}      `json:"country"`
	Name             string           `json:"name"`
	NodeID           string           `json:"node_id"`
	Region           string           `json:"region"`
	Uptime           float64          `json:"uptime"`
}

type ComputerProvider struct {
	ActiveDeployment int64       `json:"active_deployment"`
	AllowedNodes     interface{} `json:"allowed_nodes"`
	Autobid          int64       `json:"autobid"`
	City             string      `json:"city"`
	Country          interface{} `json:"country"`
	CreatedAt        string      `json:"created_at"`
	DeletedAt        interface{} `json:"deleted_at"`
	LastActiveAt     interface{} `json:"last_active_at"`
	Lat              float64     `json:"lat"`
	Lon              float64     `json:"lon"`
	Machines         []*Machine  `json:"machines"`
	MultiAddress     string      `json:"multi_address"`
	Name             string      `json:"name"`
	NodeID           string      `json:"node_id"`
	Online           bool        `json:"online"`
	Region           string      `json:"region"`
	Score            int64       `json:"score"`
	Status           string      `json:"status"`
	UpdatedAt        string      `json:"updated_at"`
}

type Machine struct {
	CreatedAt string `json:"created_at"`
	MachineID string `json:"machine_id"`
	NodeID    string `json:"node_id"`
	Specs     Specs  `json:"specs"`
	UpdatedAt string `json:"updated_at"`
}

type Specs struct {
	CPU     Stat   `json:"cpu"`
	GPU     GPU    `json:"gpu"`
	Memory  Stat   `json:"memory"`
	Model   string `json:"model"`
	Storage Stat   `json:"storage"`
	Vcpu    Stat   `json:"vcpu"`
}

type Stat struct {
	Free  string `json:"free"`
	Total string `json:"total"`
	Used  string `json:"used"`
}

type GPU struct {
	AttachedGpus  int64     `json:"attached_gpus"`
	CudaVersion   string    `json:"cuda_version"`
	Details       []*Detail `json:"details"`
	DriverVersion string    `json:"driver_version"`
}

type Detail struct {
	Bar1MemoryUsage Stat   `json:"bar1_memory_usage"`
	FbMemoryUsage   Stat   `json:"fb_memory_usage"`
	ProductName     string `json:"product_name"`
	Status          string `json:"status"`
}
