package logutil

import log "github.com/sirupsen/logrus"

func ConfigureLogger(verbose bool) {
	level := log.InfoLevel
	if verbose {
		level = log.TraceLevel
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(level)
}
