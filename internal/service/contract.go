package service

import (
	"github.com/nmluci/stellarcd/pkg/dto"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Ping() (pingResponse dto.PublicPingResponse)
}

type service struct {
	logger *logrus.Entry
	conf   *serviceConfig
}

type serviceConfig struct {
}

type NewServiceParams struct {
	Logger *logrus.Entry
}

func NewService(params *NewServiceParams) Service {
	return &service{
		logger: params.Logger,
		conf:   &serviceConfig{},
	}
}
