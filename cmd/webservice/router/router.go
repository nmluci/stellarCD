package router

import (
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/nmluci/go-backend/cmd/webservice/handler"
	"github.com/nmluci/go-backend/internal/config"
	"github.com/nmluci/go-backend/internal/service"
	"github.com/sirupsen/logrus"
)

type InitRouterParams struct {
	Logger  *logrus.Entry
	Service service.Service
	Ec      *echo.Echo
	Conf    *config.Config
}

func Init(params *InitRouterParams) {
	pprof.Register(params.Ec)

	params.Ec.GET(PingPath, handler.HandlePing(params.Service.Ping))
}
