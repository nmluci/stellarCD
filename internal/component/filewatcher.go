package component

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nmluci/stellarcd/internal/config"
	"github.com/rs/zerolog"
)

func InitFileWatcher(logger zerolog.Logger) (w *fsnotify.Watcher, err error) {
	w, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Warn().Err(err).Msg("failed to init fsnotify")
		return
	}

	conf := config.Get()
	w.Add(filepath.FromSlash(conf.DeployConfigPath))

	return
}

func WatchFilechange(logger zerolog.Logger, watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			if event.Has(fsnotify.Write) {
				logger.Info().Msg("reload config file")
				config.ReloadDeploymentConfig()
			}
		case err := <-watcher.Errors:
			logger.Warn().Err(err).Send()
		}
	}
}
