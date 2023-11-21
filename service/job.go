package service

import (
	"time"

	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/module/log"
)

func RunJobs() {
	var cjs JobService
	go cjs.Run()
}

type JobService struct {
	DBService
}

func (s *JobService) Run() {
	timer := time.NewTicker(time.Minute * 1)
	for range timer.C {
		s.SyncDeployment()
	}
}

func (s *JobService) SyncDeployment() (err error) {
	var deployments []*model.Deployment
	if err = s.DB().Model(model.Deployment{}).Where("result_url = ''").Find(&deployments).Error; err != nil {
		log.Error(err)
		return
	}
	for _, dp := range deployments {
		if err := s.LagrangeSync(dp); err != nil {
			log.Error()
		}
	}
	return
}
