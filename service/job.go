package service

import (
	"strings"
	"sync"
	"time"

	"github.com/FogMeta/libra-os/model"
	"github.com/FogMeta/libra-os/module/lagrange"
	"github.com/FogMeta/libra-os/module/log"
)

func RunJobs() {
	var cjs JobService
	go cjs.Run()
	go SyncMachines()
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

var machines sync.Map

func SyncMachines() {
	syncMachines()
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		syncMachines()
	}
}

func syncMachines() (err error) {
	data, err := lagClient.Machines()
	if err != nil {
		log.Error(err)
		return
	}
	for _, hardware := range data.Hardware {
		machine := MachineInfo{
			Hardware: *hardware,
		}
		list := strings.Split(hardware.HardwareDescription, " \u00b7 ")
		if len(list) == 3 {
			if hardware.HardwareType == "GPU" {
				machine.HardwareSpecs.GPU = list[0]
			}
			machine.HardwareSpecs.CPU = list[1]
			machine.HardwareSpecs.Memory = list[2]
		}
		machines.Store(hardware.HardwareName, machine)
	}
	return
}

func CfgMachine(name string) (info any, ok bool) {
	return machines.Load(name)
}

type MachineInfo struct {
	HardwareSpecs HardwareSpecs `json:"hardware_specs"`

	lagrange.Hardware
}

type HardwareSpecs struct {
	CPU    string `json:"cpu"`
	GPU    string `json:"gpu"`
	Memory string `json:"memory"`
}
