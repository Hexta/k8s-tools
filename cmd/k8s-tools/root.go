package main

import (
	"os"
	"path/filepath"

	"github.com/Hexta/k8s-tools/internal/logutil"
	"github.com/Hexta/k8s-tools/pkg/version"
	"k8s.io/client-go/util/homedir"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cacheDir string
	verbose  bool
)

var rootCmd = &cobra.Command{
	Use:     "k8s-tools",
	Short:   "K8s toolbox",
	Version: version.Version(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logutil.ConfigureLogger(verbose)
		err := os.MkdirAll(cacheDir, 0o755)
		if err != nil {
			log.Fatalf("error creating cache directory: %v", err)
		}

		return nil
	},
}

func init() {
	defaultCacheDir := ""
	if home := homedir.HomeDir(); home != "" {
		defaultCacheDir = filepath.Join(home, ".cache", "k8s-tools")
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().StringVar(&cacheDir, "cache-dir", defaultCacheDir, "cache directory")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
