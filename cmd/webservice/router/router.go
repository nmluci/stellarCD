package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/cmd/webservice/handler"
	"github.com/nmluci/stellarcd/internal/config"
	"github.com/nmluci/stellarcd/internal/service"
	"github.com/rs/zerolog"
)

type InitRouterParams struct {
	Logger  zerolog.Logger
	Service service.Service
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Ec.GET(PingPath, handler.HandlePing(params.Service.Ping))

	params.Ec.POST(DeploymentPath, handler.HandleDeployment(params.Service.RunDeploymentJobs))
	params.Ec.OPTIONS(DeploymentPath, handler.HandleDeployment(params.Service.RunDeploymentJobs))

	if params.Conf.Environment == config.EnvironmentDev || params.Conf.Environment == config.EnvironmentLocal {
		params.Ec.POST(SimpleDeploymentPath, handler.HandleSimpleDeployment(params.Service.RunDeploymentJobs))
		params.Ec.OPTIONS(SimpleDeploymentPath, handler.HandleSimpleDeployment(params.Service.RunDeploymentJobs))
	}
}
