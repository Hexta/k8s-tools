package main

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/k8s/fetch"
	"github.com/Hexta/k8s-tools/internal/k8s/node"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "K8s node tools",
}

var nodeCmdOpts = struct {
	LabelSelector string
}{}

var nodeUtilisationCmd = &cobra.Command{
	Use:   "utilisation",
	Short: "Analyze the node utilisation",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		clientSet := k8s.GetClientSet(getKubeconfig())

		k8sInfo := k8s.NewInfo(ctx, clientSet)
		err := k8sInfo.Fetch(fetch.Options{
			LabelSelector:        nodeCmdOpts.LabelSelector,
			RetryInitialInterval: globalOptions.k8sRetryInitialInterval,
			RetryJitterPercent:   globalOptions.k8sRetryJitterPercent,
			RetryMaxAttempts:     globalOptions.k8sRetryMaxAttempts,
			RetryMaxInterval:     globalOptions.k8sRetryMaxInterval,
		})
		if err != nil {
			log.Fatalf("Failed to fetch k8s info: %v", err)
		}

		node.FormatNodeInfo(k8sInfo.Nodes)
	},
}

func init() {
	nodeUtilisationCmd.Flags().StringVarP(&nodeCmdOpts.LabelSelector, "label-selector", "l", "", "label selector")

	nodeCmd.AddCommand(nodeUtilisationCmd)
	rootCmd.AddCommand(nodeCmd)
}
