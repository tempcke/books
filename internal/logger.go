package internal

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

const logFile = "/var/log/bookserver.log"

// Logger allows for a verbose flag, required by go-migrate
type Logger struct {
	*logrus.Logger
	verbose bool
}

// NewLogger constructs and returns a Logger
func NewLogger() *Logger {
	log := &logrus.Logger{
		Out:       os.Stdout,
		Level:     logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
	}

	options := os.O_RDWR | os.O_CREATE | os.O_APPEND
	file, err := os.OpenFile(logFile, options, 0644)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	}

	return &Logger{
		Logger: log,
	}
}

// Verbose method is required by migrate.Logger
func (log Logger) Verbose() bool {
	return log.verbose
}
