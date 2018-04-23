package logging

import (
	"github.com/rs/zerolog/log"
	"github.com/twhiston/factd/lib/common/metrics"
)

// HandleError will report an error to the logging, increase the prom error count and continue]
// Therefore this should be used to handle non fatal errors
func HandleError(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		metrics.PromErrorCount.Add(1)
	}
}

// Fatal will print the error and halt program execution, therefore this should only be used
// for critical issues, preferably only during startup time
func Fatal(err error) {
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
