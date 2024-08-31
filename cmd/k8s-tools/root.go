package main

import (
	"os"

	"github.com/Hexta/k8s-tools/internal/logutil"
	"github.com/Hexta/k8s-tools/pkg/version"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:     "k8s-tools",
	Short:   "K8s toolbox",
	Version: version.Version(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logutil.ConfigureLogger(verbose)

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
