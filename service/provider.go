package service

import (
	"sort"

	"github.com/FogMeta/libra-os/module/lagrange"
	"github.com/FogMeta/libra-os/module/log"
)

type ProviderService struct {
	DBService
}

func (s *ProviderService) ProviderList(uid int, region string) (providers []*lagrange.Provider, err error) {
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	providers, err = lagClient.WithAPIKey(user.APIKey).ProviderList(region)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *ProviderService) Provider(uid int, providerID int) (provider *lagrange.Provider, err error) {
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	provider, err = lagClient.WithAPIKey(user.APIKey).Provider(providerID)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *ProviderService) ProviderDistribution(uid int, region string) (distributions []*lagrange.Distribution, err error) {
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	distributions, err = lagClient.WithAPIKey(user.APIKey).ProviderDistribution(region)
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *ProviderService) ResourceSummary(uid int) (resource *lagrange.ProviderResource, err error) {
	user, err := s.User(uid, true)
	if err != nil {
		return
	}
	resource, err = lagClient.WithAPIKey(user.APIKey).ResourceSummary()
	if err != nil {
		log.Error(err)
	}
	return
}

func (s *ProviderService) Machines() (resource *lagrange.HardwareData, err error) {
	resource, err = lagClient.Machines()
	if err != nil {
		log.Error(err)
	}
	sort.Slice(resource.Hardware, func(i, j int) bool {
		if resource.Hardware[i].HardwareType < resource.Hardware[j].HardwareType {
			return true
		}
		if resource.Hardware[i].HardwareType == resource.Hardware[j].HardwareType &&
			resource.Hardware[i].HardwareStatus < resource.Hardware[j].HardwareStatus {
			return true
		}
		if resource.Hardware[i].HardwareType == resource.Hardware[j].HardwareType &&
			resource.Hardware[i].HardwareStatus == resource.Hardware[j].HardwareStatus &&
			resource.Hardware[i].HardwareID < resource.Hardware[j].HardwareID {
			return true
		}
		return false
	})
	return
}

func (s *ProviderService) Dashboard() (dashboard *lagrange.Dashboard, err error) {
	dashboard, err = lagClient.Dashboard()
	if err != nil {
		log.Error(err)
	}
	return
}
