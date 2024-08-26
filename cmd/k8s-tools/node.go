package main

import (
	"context"
	"path/filepath"

	"k8s-tools/internal/pkg/nodeutilisation"

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
		nodeutilisation.PrintNodeUtilisation(ctx, kubeconfig, nodepool)
	},
}

var (
	kubeconfig string
	nodepool   string
)

func init() {
	kubeconfigDefaultPath := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfigDefaultPath = filepath.Join(home, ".kube", "config")
	}

	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", kubeconfigDefaultPath, "kubeconfig file")

	nodeUtilisationCmd.Flags().StringVar(&nodepool, "nodepool", "", "node pool name")

	nodeCmd.AddCommand(nodeUtilisationCmd)
	rootCmd.AddCommand(nodeCmd)
}
