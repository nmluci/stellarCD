package service

import (
	"context"

	"github.com/nmluci/stellarcd/internal/config"
	"github.com/nmluci/stellarcd/pkg/dto"
	"github.com/nmluci/stellarcd/pkg/errs"
)

var (
	tagLoggerRunDeploymentJobs = "[RunDeploymentJobs]"
)

func (s *service) RunDeploymentJobs(ctx context.Context, payload *dto.WebhoookRequest) (err error) {
	dConf := config.GetDeploymentConfig()

	jobs, ok := dConf[payload.JobID]
	if !ok {
		s.logger.Info().Str("jobID", payload.JobID).Msg("no matching job found")
		return errs.ErrNotFound
	}
	jobs.ID = payload.JobID

	go func() {
		err = s.deployWorker.InsertJob(&jobs, payload.Webhook)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to inserting new job")
			return
		}
	}()

	return
}
