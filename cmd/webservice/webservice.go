package webservice

import (
	"github.com/labstack/echo/v4"
	"github.com/nmluci/stellarcd/cmd/webservice/router"
	"github.com/nmluci/stellarcd/internal/config"
	"github.com/nmluci/stellarcd/internal/service"
	"github.com/sirupsen/logrus"
)

const logTagStartWebservice = "[StartWebservice]"

func Start(conf *config.Config, logger *logrus.Entry) {
	// db, err := component.InitMariaDB(&component.InitMariaDBParams{
	// 	Conf:   &conf.MariaDBConfig,
	// 	Logger: logger,
	// })

	// if err != nil {
	// 	logger.Fatalf("%s initializing maria db: %+v", logTagStartWebservice, err)
	// }

	// mongo, err := component.InitMongoDB(&component.InitMongoDBParams{
	// 	Conf:   &conf.MongoDBConfig,
	// 	Logger: logger,
	// })

	// if err != nil {
	// 	logger.Fatalf("%s initializing maria db: %+v", logTagStartWebservice, err)
	// }

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

	service := service.NewService(&service.NewServiceParams{
		Logger: logger,
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

	logger.Infof("%s starting service, listening to: %s", logTagStartWebservice, conf.ServiceAddress)

	if err := ec.Start(conf.ServiceAddress); err != nil {
		logger.Errorf("%s starting service, cause: %+v", logTagStartWebservice, err)
	}
}
