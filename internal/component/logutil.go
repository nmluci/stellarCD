package component

import (
	"os"
	"runtime"
	"strings"

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

func CallerNameHook() zerolog.HookFunc {
	return func(e *zerolog.Event, l zerolog.Level, msg string) {
		pc, _, _, ok := runtime.Caller(4)
		if !ok {
			return
		}

		funcname := runtime.FuncForPC(pc).Name()
		fn := funcname[strings.LastIndex(funcname, "/")+1:]
		e.Str("caller", fn)
	}
}

func NewLogger(params NewLoggerParams) zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Str("service", params.ServiceName).Logger().Hook(CallerNameHook())
}
