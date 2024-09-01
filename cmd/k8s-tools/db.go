package main

import (
	"context"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/db"
	"github.com/Hexta/k8s-tools/internal/nodeutil"
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
		nodes := nodeutil.ListNodes(ctx, kubeconfig, labelSelector)
		err := db.InitDB(ctx, cacheDir, nodes)

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

func init() {
	dbCmd.AddCommand(initDBCmd)
	dbCmd.AddCommand(queryDBCmd)

	rootCmd.AddCommand(dbCmd)
}
