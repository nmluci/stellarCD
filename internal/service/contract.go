package service

import (
	"context"

	"github.com/nmluci/gostellar"
	"github.com/nmluci/stellarcd/internal/worker"
	"github.com/nmluci/stellarcd/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)

	NotifyError(msg string, reqID string, jobName string)
	NotifyInfo(msg string, reqID string, jobName string, versionTag string)
	RunDeploymentJobs(ctx context.Context, payload *dto.WebhoookRequest) (err error)
}

type service struct {
	logger       *logrus.Entry
	conf         *serviceConfig
	stellarRPC   *gostellar.StellarRPC
	goStellar    *gostellar.GoStellar
	deployWorker worker.DeploymentWorker
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Logger       *logrus.Entry
	StellarRPC   *gostellar.StellarRPC
	GoStellar    *gostellar.GoStellar
	DeployWorker worker.DeploymentWorker
}

func NewService(params *NewServiceParams) Service {
	return &service{
		logger:       params.Logger,
		conf:         &serviceConfig{},
		goStellar:    params.GoStellar,
		stellarRPC:   params.StellarRPC,
		deployWorker: params.DeployWorker,
	}
}
