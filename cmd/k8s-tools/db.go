package main

import (
	"fmt"

	"github.com/Hexta/k8s-tools/internal/db"
	"github.com/Hexta/k8s-tools/internal/format"
	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/k8s/fetch"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newDBCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "DB with K8s information",
	}
}

func newInitDBCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Init DB",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			restConfig, err := k8s.GetRestConfig(getKubeconfig(), globalOptions.Context)
			if err != nil {
				log.Fatalf("Failed to create rest config: %v", err)
			}

			clientSet, err := k8s.GetClientSet(restConfig)
			if err != nil {
				log.Fatalf("Failed to create clientset: %v", err)
			}

			dynamicClient, err := k8s.GetDynamicClient(restConfig)
			if err != nil {
				log.Fatalf("Failed to create dynamic client: %v", err)
			}

			apiExtClient, err := k8s.GetAPIExtensionsClient(restConfig)
			if err != nil {
				log.Fatalf("Failed to create apiextensions client: %v", err)
			}

			k8sInfo := k8s.NewInfo(ctx, clientSet, dynamicClient, apiExtClient)
			err = k8sInfo.Fetch(fetch.Options{
				RetryInitialInterval: globalOptions.K8sRetryInitialInterval,
				RetryJitterPercent:   globalOptions.K8sRetryJitterPercent,
				RetryMaxAttempts:     globalOptions.K8sRetryMaxAttempts,
				RetryMaxInterval:     globalOptions.K8sRetryMaxInterval,
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

	return cmd
}

func newQueryDBCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "query QUERY",
		Short:   "Query DB",
		Args:    cobra.ExactArgs(1),
		Example: `query "SELECT * FROM k8s.nodes LIMIT 10"`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			data, err := db.Query(ctx, getCacheDir(), args[0])
			if err != nil {
				log.Fatalf("Failed to query DB: %v", err)
			}

			output, err := format.Apply(globalOptions.Format, globalOptions.FormatOptions, data)
			if err != nil {
				log.Fatalf("Failed to format output: %v", err)
			}

			fmt.Println(output)
		},
	}
}

func newTUIDBCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "TUI for DB",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
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
