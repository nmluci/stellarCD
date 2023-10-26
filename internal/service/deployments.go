package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"

	"github.com/nmluci/stellarcd/internal/config"
	"github.com/nmluci/stellarcd/pkg/dto"
	"github.com/nmluci/stellarcd/pkg/errs"
)

func (s *service) RunDeploymentJobs(ctx context.Context, payload *dto.WebhoookRequest) (err error) {
	dConf := config.GetDeploymentConfig()

	jobs, ok := dConf[payload.JobID]
	if !ok {
		s.logger.Info().Str("jobID", payload.JobID).Msg("no matching job found")
		return errs.ErrNotFound
	}
	jobs.ID = payload.JobID

	if jobs.SignatureVal != "" { // check if jobs has static key
		if payload.HeaderMap[jobs.SignatureHeader][0] != jobs.SignatureVal {
			s.logger.Error().Err(err).Str("job-id", jobs.ID).Msg("request doesnt have matching signature")
			return errs.ErrNoAccess
		}

	} else if jobs.SignatureSecret != "" { // check if jobs has checksum
		reqHash := payload.HeaderMap[jobs.SignatureHeader][0]
		if reqHash == "" {
			s.logger.Error().Err(err).Str("job-id", jobs.ID).Msg("request signature is empty")
			return errs.ErrBadRequest
		}

		digest := hmac.New(sha256.New, []byte(jobs.SignatureSecret))
		digest.Write(payload.RawBody)
		checksum := digest.Sum(nil)

		if hmac.Equal(checksum, []byte(reqHash)) {
			s.logger.Error().Err(err).Str("job-id", jobs.ID).Msg("request doesnt have matching signature")
			return errs.ErrNoAccess
		}
	}

	go func() {
		err = s.deployWorker.InsertJob(&jobs, payload.Webhook)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to inserting new job")
			return
		}
	}()

	return
}
