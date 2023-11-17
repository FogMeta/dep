package service

import "github.com/FogMeta/libra-os/module/lagrange"

type ProviderService struct {
	DBService
}

func (s *ProviderService) ProviderList(uid int, region string) (providers []*lagrange.Provider, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	providers, err = lagClient.WithAPIKey(user.APIKey).ProviderList(region)
	return
}

func (s *ProviderService) Provider(uid int, providerID string) (provider *lagrange.Provider, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	provider, err = lagClient.WithAPIKey(user.APIKey).Provider(providerID)
	return
}

func (s *ProviderService) ProviderDistribution(uid int, region string) (distributions []*lagrange.Distribution, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	distributions, err = lagClient.WithAPIKey(user.APIKey).ProviderDistribution(region)
	return
}

func (s *ProviderService) ResourceSummary(uid int) (resource *lagrange.ProviderResource, err error) {
	user, err := s.User(uid)
	if err != nil {
		return
	}
	resource, err = lagClient.WithAPIKey(user.APIKey).ResourceSummary()
	return
}
