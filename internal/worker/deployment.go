package worker

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/nmluci/gostellar"
	"github.com/nmluci/stellarcd/internal/indto"
	"github.com/nmluci/stellarcd/pkg/errs"
	"github.com/sirupsen/logrus"
)

var (
	tagLoggerDeploymentWorker = "[DeploymentWorker]"
)

type DeploymentJob struct {
	TaskID string
	Tag    string

	Meta *indto.DeploymentJobs
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

	if job.TriggerRegex != "" {
		re, err := regexp.Compile(job.TriggerRegex)
		if err != nil {
			dw.logger.Errorf("%s failed to validate regex matching err: %+v", tagLoggerDeploymentWorker, err)
			dw.NotifyError("failed to validate regex matching", task.TaskID, task.Meta.ID)
			return errs.ErrBadRequest
		}

		_, ok := payload[job.TriggerKey].(string)
		if !ok {
			dw.NotifyError("failed to find trigger", task.TaskID, task.Meta.ID)
			return errs.ErrNotFound
		}

		if tag := re.FindString(payload[job.TriggerKey].(string)); tag == "" {
			dw.NotifyError("failed to find tag", task.TaskID, task.Meta.ID)
			return errs.ErrNotFound
		} else {
			task.Tag = tag
		}
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
			continue
		}

		cmd := exec.Command(cmdPath)
		cmd.Dir = job.Meta.WorkingDir
		cmd.Args = []string{job.Meta.Command}

		msg, err := cmd.CombinedOutput()
		if err != nil {
			dw.logger.Errorf("%s command err: %+v", tagLoggerDeploymentWorker, err)
			dw.NotifyError(err.Error(), job.TaskID, job.Meta.ID)
			continue
		}

		dw.logger.Infof("%s", msg)
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
