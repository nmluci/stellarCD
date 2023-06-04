package component

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

type NewLoggerParams struct {
	PrettyPrint bool
	ServiceName string
}

type UTCFormatter struct {
	logrus.Formatter
}

func NewLogger(params NewLoggerParams) zerolog.Logger {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return log
}
