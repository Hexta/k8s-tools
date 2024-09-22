package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Hexta/k8s-tools/internal/db"
	"github.com/Hexta/k8s-tools/internal/k8s"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type dbCmdOptions struct {
	k8sRetryInitialInterval time.Duration
	k8sRetryJitterPercent   uint64
	k8sRetryMaxAttempts     uint64
	k8sRetryMaxInterval     time.Duration
}

func newDBCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "DB with K8s information",
	}
}

func newInitDBCmd() *cobra.Command {
	cmdOptions := dbCmdOptions{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init DB",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			clientSet := k8s.GetClientSet(getKubeconfig())

			k8sInfo := k8s.NewInfo(ctx, clientSet)
			err := k8sInfo.Fetch(k8s.FetchOptions{
				RetryInitialInterval: cmdOptions.k8sRetryInitialInterval,
				RetryJitterPercent:   cmdOptions.k8sRetryJitterPercent,
				RetryMaxAttempts:     cmdOptions.k8sRetryMaxAttempts,
				RetryMaxInterval:     cmdOptions.k8sRetryMaxInterval,
			})
			if err != nil {
				log.Fatalf("Failed to fetch K8s info: %v", err)
			}

			err = db.InitDB(ctx, getCacheDir(), k8sInfo)
			if err != nil {
				log.Fatalf("Failed to init DB: %v", err)
			}
		},
	}

	cmd.Flags().DurationVarP(&cmdOptions.k8sRetryInitialInterval, "k8s-retry-initial-interval", "", time.Second, "Initial interval for Kubernetes API retry")
	cmd.Flags().Uint64VarP(&cmdOptions.k8sRetryJitterPercent, "k8s-retry-jitter-percent", "", 50, "Jitter percent for Kubernetes API retry")
	cmd.Flags().Uint64VarP(&cmdOptions.k8sRetryMaxAttempts, "k8s-retry-max-attempts", "", 5, "Maximum number of attempts for Kubernetes API retry")
	cmd.Flags().DurationVarP(&cmdOptions.k8sRetryMaxInterval, "k8s-retry-max-interval", "", 10*time.Second, "Maximum interval between retries for Kubernetes API")
	return cmd
}

func newQueryDBCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Query DB",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			result, err := db.Query(ctx, getCacheDir(), args[0])

			if err != nil {
				log.Fatalf("Failed to query DB: %v", err)
			}

			fmt.Print(result)
		},
	}
}

func newTUIDBCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "TUI for DB",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			db.RunTUI(ctx, getCacheDir())
		},
	}
}

func init() {
	dbCmd := newDBCmd()
	dbCmd.AddCommand(newInitDBCmd())
	dbCmd.AddCommand(newTUIDBCmd())
	dbCmd.AddCommand(newQueryDBCmd())

	rootCmd.AddCommand(dbCmd)
}
