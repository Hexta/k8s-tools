package main

import (
	"context"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/containerutil"
	"github.com/Hexta/k8s-tools/internal/db"
	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/nodeutil"
	"github.com/Hexta/k8s-tools/internal/podutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "DB with K8s information",
}

var initDBCmd = &cobra.Command{
	Use:   "init",
	Short: "Init DB",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		clientSet := k8s.GetClientSet(kubeconfig)
		nodes := nodeutil.ListNodes(ctx, kubeconfig, labelSelector)
		pods := podutil.ListPods(ctx, clientSet, labelSelector)
		containers := containerutil.ListContainers(ctx, clientSet, labelSelector)
		err := db.InitDB(ctx, cacheDir, nodes, pods, containers)

		if err != nil {
			log.Fatalf("Failed to init DB: %v", err)
		}
	},
}

var queryDBCmd = &cobra.Command{
	Use:   "query",
	Short: "Query DB",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		result, err := db.Query(ctx, cacheDir, args[0])

		if err != nil {
			log.Fatalf("Failed to query DB: %v", err)
		}

		fmt.Print(result)
	},
}

var tuiDBCmd = &cobra.Command{
	Use:   "tui",
	Short: "TUI for DB",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		db.RunTUI(ctx, cacheDir)

	},
}

func init() {
	dbCmd.AddCommand(initDBCmd)
	dbCmd.AddCommand(tuiDBCmd)
	dbCmd.AddCommand(queryDBCmd)

	rootCmd.AddCommand(dbCmd)
}
