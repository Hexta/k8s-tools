package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Hexta/k8s-tools/internal/format"
	"github.com/Hexta/k8s-tools/internal/logutil"
	"github.com/Hexta/k8s-tools/pkg/version"
	"k8s.io/client-go/util/homedir"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	globalOptions = struct {
		CacheDir                string
		Format                  format.Format
		FormatOptions           format.Options
		Kubeconfig              string
		Verbose                 bool
		K8sRetryInitialInterval time.Duration
		K8sRetryJitterPercent   uint64
		K8sRetryMaxAttempts     uint64
		K8sRetryMaxInterval     time.Duration
	}{
		Format: format.TableFormat,
	}
)

var rootCmd = &cobra.Command{
	Use:     "k8s-tools",
	Short:   "K8s toolbox",
	Version: version.Version(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logutil.ConfigureLogger(globalOptions.Verbose)

		return ensureCacheDir()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalOptions.Verbose, "verbose", "v", false, "verbose")
	rootCmd.PersistentFlags().StringVar(&globalOptions.CacheDir, "cache-dir", "", "cache directory")
	rootCmd.PersistentFlags().StringVar(&globalOptions.Kubeconfig, "kubeconfig", "", "kubeconfig file")

	rootCmd.PersistentFlags().DurationVarP(&globalOptions.K8sRetryInitialInterval, "k8s-retry-initial-interval", "", time.Second, "Initial interval for Kubernetes API retry")
	rootCmd.PersistentFlags().Uint64VarP(&globalOptions.K8sRetryJitterPercent, "k8s-retry-jitter-percent", "", 50, "Jitter percent for Kubernetes API retry")
	rootCmd.PersistentFlags().Uint64VarP(&globalOptions.K8sRetryMaxAttempts, "k8s-retry-max-attempts", "", 5, "Maximum number of attempts for Kubernetes API retry")
	rootCmd.PersistentFlags().DurationVarP(&globalOptions.K8sRetryMaxInterval, "k8s-retry-max-interval", "", 10*time.Second, "Maximum interval between retries for Kubernetes API")

	rootCmd.PersistentFlags().VarP(&globalOptions.Format, "output", "o", "output format (json, table, vertical)")
	rootCmd.PersistentFlags().BoolVar(&globalOptions.FormatOptions.NoHeaders, "no-headers", false, "do not print headers")

	err := rootCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			format.JSONFormat, format.TableFormat, format.VerticalFormat,
		}, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		log.Fatalf("error registering flag completion: %v", err)
	}
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

// Ensure the cache directory exists and is writable.
func ensureCacheDir() error {
	cacheDir := getCacheDir()
	if cacheDir == "" {
		return fmt.Errorf("unable to determine cache directory")
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("error creating cache directory: %w", err)
	}
	return nil
}
