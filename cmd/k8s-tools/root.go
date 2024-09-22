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
	globalOptions = struct {
		Verbose    bool
		CacheDir   string
		Kubeconfig string
	}{}
)

var rootCmd = &cobra.Command{
	Use:     "k8s-tools",
	Short:   "K8s toolbox",
	Version: version.Version(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logutil.ConfigureLogger(globalOptions.Verbose)
		err := os.MkdirAll(getCacheDir(), 0o755)
		if err != nil {
			log.Fatalf("error creating cache directory: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalOptions.Verbose, "verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().StringVar(&globalOptions.CacheDir, "cache-dir", "", "cache directory")
	rootCmd.PersistentFlags().StringVar(&globalOptions.Kubeconfig, "kubeconfig", "", "kubeconfig file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func getCacheDir() string {
	defaultCacheDir := ""
	if home := homedir.HomeDir(); home != "" {
		defaultCacheDir = filepath.Join(home, ".cache", "k8s-tools")
	}

	if globalOptions.CacheDir != "" {
		return globalOptions.CacheDir
	}

	return defaultCacheDir
}

func getKubeconfig() string {
	defaultPath := ""
	if home := homedir.HomeDir(); home != "" {
		defaultPath = filepath.Join(home, ".kube", "config")
	}

	if globalOptions.Kubeconfig != "" {
		return globalOptions.Kubeconfig
	}

	return defaultPath
}
