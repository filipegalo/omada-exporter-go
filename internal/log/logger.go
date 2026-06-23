package Log

import (
	"os"

	"github.com/rs/zerolog"

	"omada_exporter_go/internal"
)

const (
	defaultLogLevel    = zerolog.InfoLevel
	callerFramesToSkip = 3
)

var (
	logger zerolog.Logger
)

func Init() {
	conf := internal.GetConfig()
	logLevel, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		logLevel = defaultLogLevel
	}

	logger = zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		CallerWithSkipFrameCount(callerFramesToSkip).
		Logger()
}
