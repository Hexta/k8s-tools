package main

import (
	"context"
	"path/filepath"

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
		nodes := nodeutil.ListNodes(ctx, kubeconfig, labelSelector)
		nodeutil.FormatNodeInfo(nodes)
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
