package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a wrapper around logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set log level based on environment variable or default to info
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logrusLogger.SetLevel(logrus.DebugLevel)
	case "warn":
		logrusLogger.SetLevel(logrus.WarnLevel)
	case "error":
		logrusLogger.SetLevel(logrus.ErrorLevel)
	default:
		logrusLogger.SetLevel(logrus.InfoLevel)
	}

	return &Logger{Logger: logrusLogger}
}