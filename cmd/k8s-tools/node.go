package main

import (
	"context"
	"path/filepath"

	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/nodeutil"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
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
		clientSet := k8s.GetClientSet(kubeconfig)

		k8sInfo := k8s.NewInfo(ctx, clientSet)
		k8sInfo.Fetch(labelSelector, labelSelector)

		nodeutil.FormatNodeInfo(k8sInfo.Nodes)
	},
}

var (
	kubeconfig    string
	labelSelector string
)

func init() {
	kubeconfigDefaultPath := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultPath = filepath.Join(home, ".kube", "config")
	}

	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", kubeconfigDefaultPath, "kubeconfig file")

	nodeUtilisationCmd.Flags().StringVarP(&labelSelector, "label-selector", "l", "", "label selector")

	nodeCmd.AddCommand(nodeUtilisationCmd)
	rootCmd.AddCommand(nodeCmd)
}
