package main

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/k8s/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "K8s node tools",
}

var nodeUtilisationCmd = &cobra.Command{
	Use:   "utilisation",
	Short: "Analyze the node utilisation",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		clientSet := k8s.GetClientSet(getCacheDir())

		k8sInfo := k8s.NewInfo(ctx, clientSet)
		err := k8sInfo.Fetch(k8s.FetchOptions{})
		if err != nil {
			log.Fatalf("Failed to fetch k8s info: %v", err)
		}

		node.FormatNodeInfo(k8sInfo.Nodes)
	},
}

var (
	labelSelector string
)

func init() {
	nodeUtilisationCmd.Flags().StringVarP(&labelSelector, "label-selector", "l", "", "label selector")

	nodeCmd.AddCommand(nodeUtilisationCmd)
	rootCmd.AddCommand(nodeCmd)
}
