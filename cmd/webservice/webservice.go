package webservice

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/gostellar"
	"github.com/nmluci/stellarcd/cmd/webservice/router"
	"github.com/nmluci/stellarcd/internal/component"
	"github.com/nmluci/stellarcd/internal/config"
	"github.com/nmluci/stellarcd/internal/service"
	"github.com/nmluci/stellarcd/internal/worker"
	"github.com/rs/zerolog"
)

const logTagStartWebservice = "[StartWebservice]"

func Start(conf *config.Config, logger zerolog.Logger) {
	// redis, err := component.InitRedis(&component.InitRedisParams{
	// 	Conf:   &conf.RedisConfig,
	// 	Logger: logger,
	// })

	// if err != nil {
	// 	logger.Fatalf("%s initalizing redis: %+v", logTagStartWebservice, err)
	// }

	ec := echo.New()
	ec.HideBanner = true
	ec.HidePort = true

	gs := gostellar.NewGoStellar(gostellar.NewGoStellarParams{
		Logger:      &logger,
		ServiceName: conf.ServiceID,
	})

	worker := worker.NewDeploymentWorker(&worker.NewDeploymentWorkerParams{
		Logger:    logger,
		GoStellar: &gs,
	})

	service := service.NewService(&service.NewServiceParams{
		Logger:       logger,
		StellarRPC:   gs.StellarRPC,
		GoStellar:    &gs,
		DeployWorker: worker,
	})

	// psWorker := pubsub.NewFileSub(pubsub.NewFilePubSubParams{
	// 	Logger: logger,
	// 	// Redis:   redis,
	// 	Service: service,
	// })

	router.Init(&router.InitRouterParams{
		Logger:  logger,
		Service: service,
		Ec:      ec,
		Conf:    conf,
	})

	config.ReloadDeploymentConfig()

	go worker.StartWorker()

	watcher, err := component.InitFileWatcher(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init filewatcher")
	}
	go component.WatchFilechange(logger, watcher)

	logger.Info().Msgf("starting service, listening to: %s", conf.ServiceAddress)
	if err := ec.Start(conf.ServiceAddress); err != nil {
		logger.Error().Err(err).Msg("failed to start service")
	}

	watcher.Close()
	worker.StopWorker()
}
