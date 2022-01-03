package logger

import (
	"github.com/nbs-go/nlogger"
	"os"
)

func init() {
	// Get log level from env
	logLevelStr, _ := os.LookupEnv(nlogger.EnvLogLevel)
	logLevel := nlogger.ParseLevel(logLevelStr)

	// Get namespace
	namespace, _ := os.LookupEnv("APP_NAME")

	// Register log
	nlogger.Register(NewJson(logLevel, nil, namespace))
}
