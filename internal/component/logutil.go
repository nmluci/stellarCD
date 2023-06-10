package component

import (
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/nmluci/stellarcd/internal/config"
	"github.com/rs/zerolog"
)

type NewLoggerParams struct {
	PrettyPrint bool
	ServiceName string
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
	var output io.Writer

	if env := os.Getenv("ENVIRONMENT"); env == string(config.EnvironmentLocal) {
		output = zerolog.ConsoleWriter{
			Out: os.Stdout,
		}
	} else {
		output = os.Stdout
	}

	return zerolog.New(output).With().Timestamp().Str("service", params.ServiceName).Logger().Hook(CallerNameHook())
}
