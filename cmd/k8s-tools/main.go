package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	Execute()
}

func configureLogger(verbose bool) {
	level := log.InfoLevel
	if verbose {
		level = log.TraceLevel
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetLevel(level)
}
