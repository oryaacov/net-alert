package logging

import (
	"fmt"
	"io"
	"net-alert/pkg/config"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

//InitLog init logger using logrus
func InitLog(config *config.Configuration) {
	mw := io.MultiWriter(getOutputStreams(config)...)
	log.SetOutput(mw)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	if lvl, err := log.ParseLevel(config.Log.LogLevel); err != nil {
		panic(err)
	} else {
		// Only log the warning severity or above.
		log.SetLevel(lvl)
	}
}

//LogInfo logs info using logrus
func LogInfo(args ...interface{}) {
	log.Info(args...)
}

//LogWarn logs warning using logrus
func LogWarn(args ...interface{}) {
	log.Warn(args...)
}

//LogError logs error using logrus
func LogError(args ...interface{}) {
	log.Error(args...)
}

//LogFatal logs fatal error using logrus
func LogFatal(args ...interface{}) {
	log.Fatal(args...)
}

//LogPanic logs panic using logrus
func LogPanic(args ...interface{}) {
	log.Panic(args...)
}

//LogDebug logs debug info using logrus
func LogDebug(args ...interface{}) {
	log.Debug(args...)
}

func getOutputStreams(config *config.Configuration) []io.Writer {
	result := make([]io.Writer, 0)
	if config.Log.LogToConsole {
		result = append(result, os.Stdout)
	}
	if logFile, err := os.Create(fmt.Sprintf(config.Log.LogFilePath, time.Now().String()[:19])); err != nil {
		panic(err)
	} else {
		result = append(result, logFile)
	}
	return result
}
