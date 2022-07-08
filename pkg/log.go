package pkg

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	// set default formatter
	logrus.SetFormatter(&prefixed.TextFormatter{})
}

func NewScopedLogger(scope string) *logrus.Entry {
	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	}
	return logger.WithFields(logrus.Fields{"scope": scope})
}
