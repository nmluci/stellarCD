package worker

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nmluci/gostellar"
	"github.com/nmluci/gostellar/pkg/dto"
	"github.com/nmluci/stellarcd/internal/indto"
	"github.com/nmluci/stellarcd/pkg/errs"
	"github.com/rs/zerolog"
)

var (
	tagLoggerDeploymentWorker = "[DeploymentWorker]"
)

type DeploymentJob struct {
	TaskID    string
	Tag       string
	CommitMsg string

	Meta        *indto.DeploymentJobs
	WebhookCred *dto.DiscordWebhoookCred
}

type DeploymentWorker interface {
	InsertJob(job *indto.DeploymentJobs, payload map[string]interface{}) (err error)
	Executor(id int)
	StartWorker()
	StopWorker()
}

type deploymentWorker struct {
	wg        *sync.WaitGroup
	logger    zerolog.Logger
	goStellar *gostellar.GoStellar
	jobQueue  chan DeploymentJob
}

type NewDeploymentWorkerParams struct {
	Logger    zerolog.Logger
	GoStellar *gostellar.GoStellar
}

func NewDeploymentWorker(params *NewDeploymentWorkerParams) (dw DeploymentWorker) {
	dw = &deploymentWorker{
		wg:        &sync.WaitGroup{},
		logger:    params.Logger.With().Str("module", "DeploymentWorker").Logger(),
		jobQueue:  make(chan DeploymentJob, 10),
		goStellar: params.GoStellar,
	}

	return
}

func (dw *deploymentWorker) InsertJob(job *indto.DeploymentJobs, payload map[string]interface{}) (err error) {
	task := DeploymentJob{
		TaskID: uuid.NewString(),
		Meta:   job,
	}

	if job.WebhookID != "" && job.WebhookToken != "" {
		task.WebhookCred = &dto.DiscordWebhoookCred{
			WebhookID:    job.WebhookID,
			WebhookToken: job.WebhookToken,
		}
	}

	if job.TriggerRegex != "" {
		re, err := regexp.Compile(job.TriggerRegex)
		if err != nil {
			dw.logger.Error().Err(err).Msg("failed to validate regex matching")
			dw.NotifyError(task.WebhookCred, "failed to validate regex matching", task.TaskID, task.Meta.ID)
			return errs.ErrBadRequest
		}

		_, ok := payload[job.TriggerKey].(string)
		if !ok {
			dw.NotifyError(task.WebhookCred, "failed to find trigger", task.TaskID, task.Meta.ID)
			return errs.ErrNotFound
		}

		if tag := re.FindString(payload[job.TriggerKey].(string)); tag == "" {
			return errs.ErrNotFound
		} else {
			task.Tag = re.FindStringSubmatch(payload[job.TriggerKey].(string))[1]
		}

		// // TODO: Refactor nested attribute fetch
		// if msg, ok := payload["head_commit"].(map[string]interface{})["message"]; ok {
		// 	task.CommitMsg = msg
		// }
	}

	// TODO: Add SHA validation

	dw.jobQueue <- task
	dw.logger.Info().Msg("succesfully insert new job")

	return
}

func (dw *deploymentWorker) Executor(id int) {
	dw.logger.Info().Int("id", id).Msg("initialized DeploymentWorker")

	for job := range dw.jobQueue {
		dw.logger.Info().Str("jobID", job.TaskID).Msg("running job")

		var lookpath string
		if filepath.IsAbs(job.Meta.Command) || job.Meta.WorkingDir != "" {
			lookpath = job.Meta.Command
		} else {
			lookpath = filepath.Join(job.Meta.WorkingDir, job.Meta.Command)
		}

		cmdPath, err := exec.LookPath(lookpath)
		if err != nil {
			dw.logger.Error().Err(err).Send()
			dw.NotifyError(job.WebhookCred, fmt.Sprintf("lookpath err: %+v", err), job.TaskID, job.Meta.ID)
			continue
		}

		cmd := exec.Command(cmdPath)
		cmd.Dir = job.Meta.WorkingDir
		cmd.Args = []string{job.Meta.Command}
		cmd.Env = append(os.Environ(), fmt.Sprintf("BUILD_TAG=%s", job.Tag), fmt.Sprintf("BUILD_TIMESTAMP=%s", time.Now().Format("2006-01-02 15:04:05")), "BUILDKIT_PROGRESS=plain")

		cmdOut, err := cmd.StdoutPipe()
		if err != nil {
			dw.logger.Warn().Err(err).Msg("stdout pipe")
			dw.NotifyError(job.WebhookCred, fmt.Sprintf("stdout pipe err: %+v", err), job.TaskID, job.Meta.ID)
			continue
		}

		cmdScanner := bufio.NewScanner(cmdOut)
		go func() {
			for cmdScanner.Scan() {
				dw.logger.Info().Str("job", job.Meta.ID).Str("uuid", job.TaskID).Msg(cmdScanner.Text())
			}
		}()

		err = cmd.Start()
		if err != nil {
			dw.logger.Error().Err(err).Msg("failed to start deployment")
			dw.NotifyError(job.WebhookCred, err.Error(), job.TaskID, job.Meta.ID)
			continue
		}

		err = cmd.Wait()
		if err != nil {
			dw.logger.Error().Err(err).Send()
			dw.NotifyError(job.WebhookCred, err.Error(), job.TaskID, job.Meta.ID)
			continue
		}

		dw.NotifyInfo(job.WebhookCred, "deploy success", job.TaskID, job.Meta.ID, job.Tag, job.CommitMsg)
		dw.logger.Info().Str("taskID", job.TaskID).Str("jobID", job.Meta.ID).Str("tag", job.Tag).Msg("deploy success")
	}

}

func (dw *deploymentWorker) StartWorker() {
	go dw.Executor(1)
}

func (dw *deploymentWorker) StopWorker() {
	dw.wg.Wait()
	dw.logger.Warn().Msg("gracefully shutting down worker")
	close(dw.jobQueue)
}
