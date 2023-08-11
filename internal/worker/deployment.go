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
	TaskID          string
	Tag             string
	CommitMessage   string
	CommitURL       string
	CommitTimestamp string
	CommitAuthor    string

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
			dw.NotifyError(task.WebhookCred, NotifyErrorParams{
				Message: "failed to validate regex matching",
				ReqID:   task.TaskID,
				JobName: task.Meta.ID,
			})
			return errs.ErrBadRequest
		}

		_, ok := payload[job.TriggerKey].(string)
		if !ok {
			dw.NotifyError(task.WebhookCred, NotifyErrorParams{
				Message: "failed to find trigger",
				ReqID:   task.TaskID,
				JobName: task.Meta.ID,
			})
			return errs.ErrNotFound
		}

		if tag := re.FindString(payload[job.TriggerKey].(string)); tag == "" {
			return errs.ErrNotFound
		} else {
			task.Tag = re.FindStringSubmatch(payload[job.TriggerKey].(string))[1]
		}

		if msg, ok := payload["head_commit"].(map[string]interface{}); ok {
			task.CommitMessage = msg["message"].(string)
			task.CommitURL = msg["url"].(string)
			task.CommitTimestamp = msg["timestamp"].(string)
			commitAuthor := msg["author"].(map[string]interface{})
			task.CommitAuthor = fmt.Sprintf("%s <%s>", commitAuthor["username"].(string), commitAuthor["name"].(string))
		}
	}

	// TODO: Add SHA validation

	dw.jobQueue <- task
	dw.logger.Info().Msg("succesfully insert new job")

	return
}

func (dw *deploymentWorker) Executor(id int) {
	dw.logger.Info().Int("id", id).Msg("initialized DeploymentWorker")

	for job := range dw.jobQueue {
		timeStart := time.Now()
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
			dw.NotifyError(job.WebhookCred, NotifyErrorParams{
				Message: fmt.Sprintf("lookpath err: %+v", err),
				ReqID:   job.TaskID,
				JobName: job.Meta.ID,
			})
			continue
		}

		cmd := exec.Command(cmdPath)
		cmd.Dir = job.Meta.WorkingDir
		cmd.Args = []string{job.Meta.Command}
		cmd.Env = append(os.Environ(), fmt.Sprintf("BUILD_TAG=%s", job.Tag), fmt.Sprintf("BUILD_TIMESTAMP=%s", time.Now().Format("2006-01-02 15:04:05")), "BUILDKIT_PROGRESS=plain")

		cmdOut, err := cmd.StdoutPipe()
		if err != nil {
			dw.logger.Warn().Err(err).Msg("stdout pipe")
			dw.NotifyError(job.WebhookCred, NotifyErrorParams{
				Message: fmt.Sprintf("stdout pipe err: %+v", err),
				ReqID:   job.TaskID,
				JobName: job.Meta.ID,
			})
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
			dw.NotifyError(job.WebhookCred, NotifyErrorParams{
				Message: err.Error(),
				ReqID:   job.TaskID,
				JobName: job.Meta.ID,
			})
			continue
		}

		// ignoring err from cmd.Wait()
		// due weird return code
		cmd.Wait()

		dw.NotifyInfo(job.WebhookCred, NotifyInfoParams{
			Message:         "deploy success",
			ReqID:           job.TaskID,
			JobName:         job.Meta.ID,
			VersionTag:      job.Tag,
			CommitMessage:   job.CommitMessage,
			CommitURL:       job.CommitURL,
			CommitTimestamp: job.CommitTimestamp,
			CommitAuthor:    job.CommitAuthor,
			BuildTime:       fmt.Sprintf("%.2f", time.Since(timeStart).Seconds()),
		})
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
