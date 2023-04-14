package service

import (
	"context"

	"github.com/nmluci/go-backend/internal/component"
	"github.com/nmluci/go-backend/internal/repository"
	"github.com/nmluci/go-backend/internal/worker"
	"github.com/nmluci/go-backend/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
	AuthenticateSession(ctx context.Context, token string) (access context.Context, err error)
	AuthenticateService(ctx context.Context, name string) (access context.Context, err error)
}

type service struct {
	logger     *logrus.Entry
	conf       *serviceConfig
	repository repository.Repository
	stellarRPC *component.StellarRPCService
	fileWorker *worker.WorkerManager
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Logger     *logrus.Entry
	Repository repository.Repository
	StellarRPC *component.StellarRPCService
	FileWorker *worker.WorkerManager
}

func NewService(params *NewServiceParams) Service {
	return &service{
		logger:     params.Logger,
		conf:       &serviceConfig{},
		repository: params.Repository,
		stellarRPC: params.StellarRPC,
		fileWorker: params.FileWorker,
	}
}
