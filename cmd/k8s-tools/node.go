package main

import (
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
		ctx := cmd.Context()
		clientSet, err := k8s.GetClientSet(getKubeconfig(), globalOptions.Context)
		if err != nil {
			log.Fatalf("Failed to create clientset: %v", err)
		}

		k8sInfo := k8s.NewInfo(ctx, clientSet)
		err = k8sInfo.Fetch(fetch.Options{
			LabelSelector:        nodeCmdOpts.LabelSelector,
			RetryInitialInterval: globalOptions.K8sRetryInitialInterval,
			RetryJitterPercent:   globalOptions.K8sRetryJitterPercent,
			RetryMaxAttempts:     globalOptions.K8sRetryMaxAttempts,
			RetryMaxInterval:     globalOptions.K8sRetryMaxInterval,
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
