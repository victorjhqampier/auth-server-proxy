package handler

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

type LoggerHandler struct {
	logger *logrus.Logger
}

func NewLoggerHandler() *LoggerHandler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.DebugLevel)

	return &LoggerHandler{
		logger: logger,
	}
}

func (lh *LoggerHandler) Info(message string) {
	lh.logger.Info(message)
}

func (lh *LoggerHandler) Error(err error, contextMessage string) {
	lh.logger.WithError(err).
		WithField("stack", "> "+string(debug.Stack())).
		Error(contextMessage)
}

func (lh *LoggerHandler) Warning(message string) {
	lh.logger.Warn(message)
}
