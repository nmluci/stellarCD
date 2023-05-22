package worker

import (
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
	"github.com/sirupsen/logrus"
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
	logger    *logrus.Entry
	goStellar *gostellar.GoStellar
	jobQueue  chan DeploymentJob
}

type NewDeploymentWorkerParams struct {
	Logger    *logrus.Entry
	GoStellar *gostellar.GoStellar
}

func NewDeploymentWorker(params *NewDeploymentWorkerParams) (dw DeploymentWorker) {
	dw = &deploymentWorker{
		wg:        &sync.WaitGroup{},
		logger:    params.Logger,
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
			dw.logger.Errorf("%s failed to validate regex matching err: %+v", tagLoggerDeploymentWorker, err)
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
	dw.logger.Infof("%s succesfully insert new job", tagLoggerDeploymentWorker)

	return
}

func (dw *deploymentWorker) Executor(id int) {
	dw.logger.Infof("%s initialized DeploymentWorker id: %d", tagLoggerDeploymentWorker, id)

	for job := range dw.jobQueue {
		dw.logger.Infof("%s running job id: %s", tagLoggerDeploymentWorker, job.TaskID)

		var lookpath string
		if filepath.IsAbs(job.Meta.Command) || job.Meta.WorkingDir != "" {
			lookpath = job.Meta.Command
		} else {
			lookpath = filepath.Join(job.Meta.WorkingDir, job.Meta.Command)
		}

		cmdPath, err := exec.LookPath(lookpath)
		if err != nil {
			dw.logger.Printf("%s lookpath err: %+v", tagLoggerDeploymentWorker, err)
			dw.NotifyError(job.WebhookCred, fmt.Sprintf("lookpath err: %+v", err), job.TaskID, job.Meta.ID)
			continue
		}

		cmd := exec.Command(cmdPath)
		cmd.Dir = job.Meta.WorkingDir
		cmd.Args = []string{job.Meta.Command}
		cmd.Env = append(os.Environ(), fmt.Sprintf("BUILD_TAG=%s", job.Tag), fmt.Sprintf("BUILD_TIMESTAMP=%s", time.Now().Format("2006-01-02 15:04:05")))

		msg, err := cmd.CombinedOutput()
		if err != nil {
			dw.logger.Errorf("%s command err: %+v", tagLoggerDeploymentWorker, err)
			dw.NotifyError(job.WebhookCred, err.Error(), job.TaskID, job.Meta.ID)
			dw.logger.Infof("Err: %s", msg)
			continue
		}

		dw.NotifyInfo(job.WebhookCred, "deploy success", job.TaskID, job.Meta.ID, job.Tag, job.CommitMsg)
		dw.logger.Infof("%s deploy success taskID: %s, jobID: %s, tag: %s", tagLoggerDeploymentWorker, job.TaskID, job.Meta.ID, job.Tag)
	}

}

func (dw *deploymentWorker) StartWorker() {
	go dw.Executor(1)
}

func (dw *deploymentWorker) StopWorker() {
	dw.wg.Wait()
	dw.logger.Errorf("%s gracefully shutting down worker", tagLoggerDeploymentWorker)
	close(dw.jobQueue)
}
