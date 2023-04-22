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
		s.logger.Infof("%s no matching job found. (got: %s)", tagLoggerRunDeploymentJobs, payload.JobID)
		return errs.ErrNotFound
	}
	jobs.ID = payload.JobID

	go func() {
		err = s.deployWorker.InsertJob(&jobs, payload.Webhook)
		if err != nil {
			s.logger.Errorf("%s err while inserting new job: %+v", tagLoggerRunDeploymentJobs, err)
			return
		}
	}()

	return
}
