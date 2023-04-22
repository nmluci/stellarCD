package config

import (
	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
	"github.com/nmluci/stellarcd/internal/indto"
)

var deploymentJobs map[string]indto.DeploymentJobs

func ReloadDeploymentConfig() {
	temp := make(map[string]indto.DeploymentJobs)

	_, err := toml.DecodeFile("conf/deployment.toml", &temp)
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	deploymentJobs = temp
}

func GetDeploymentConfig() map[string]indto.DeploymentJobs {
	return deploymentJobs
}