package lagrange

import "net/url"

const (
	methodProviderList         = "/cp_list"
	methodProviderDistribution = "/cp_distribution"
	methodProviderDetail       = "/cp_detail"
	methodProviderAvailable    = "/cp_available"
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
