package main

import (
	"github.com/nmluci/stellarcd/cmd/webservice"
	"github.com/nmluci/stellarcd/internal/component"
	"github.com/nmluci/stellarcd/internal/config"
)

func main() {
	config.Init()
	conf := config.Get()

	logger := component.NewLogger(component.NewLoggerParams{
		ServiceName: conf.ServiceName,
		PrettyPrint: true,
	})

	webservice.Start(conf, logger)
}
