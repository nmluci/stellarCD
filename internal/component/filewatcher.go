package component

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nmluci/stellarcd/internal/config"
	"github.com/sirupsen/logrus"
)

func InitFileWatcher(logger *logrus.Entry) (w *fsnotify.Watcher, err error) {
	w, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Infof("[WatchFileChange] failed to init fsnotify err: %+v", err)
		return
	}

	conf := config.Get()
	w.Add(filepath.FromSlash(conf.DeployConfigPath))

	return
}

func WatchFilechange(logger *logrus.Entry, watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			if event.Has(fsnotify.Write) {
				logger.Infof("[WatchFileChange] reload config file")
				config.ReloadDeploymentConfig()
			}
		case err := <-watcher.Errors:
			logger.Infof("[WatchFileChange] watcher error: %+v", err)
		}
	}
}
