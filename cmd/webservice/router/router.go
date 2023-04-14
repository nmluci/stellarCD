package router

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/cmd/webservice/handler"
	"github.com/nmluci/stellarcd/internal/config"
	"github.com/nmluci/stellarcd/internal/service"
	"github.com/sirupsen/logrus"
)

type InitRouterParams struct {
	Logger  *logrus.Entry
	Service service.Service
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	params.Ec.GET(PingPath, handler.HandlePing(params.Service.Ping))
	params.Ec.POST(ReflectorPath, handler.HandleReflector())
	params.Ec.OPTIONS(ReflectorPath, handler.HandleReflector())
}
