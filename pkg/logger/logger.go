package logger

import (
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
)

func Init(debug bool) {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugf("Debug mode is enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
