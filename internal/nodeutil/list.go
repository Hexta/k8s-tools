package nodeutil

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func ListNodes(ctx context.Context, kubeconfigFile string, labelSelector string) []NodeInfo {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		logrus.Fatalf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Failed to create client: %v", err)
	}

	list, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{
		TypeMeta:      v1.TypeMeta{},
		LabelSelector: labelSelector,
	})
	if err != nil {
		logrus.Fatalf("Failed to list nodes: %v", err)
	}

	nodes := make([]NodeInfo, 0, len(list.Items))

	logrus.Debugf("Listing Pods per node - start")
	podsPerNode, err := listPodsPerNode(ctx, clientset)
	logrus.Debugf("Listing Pods per node - done. Size: %v", len(podsPerNode))
	if err != nil {
		logrus.Errorf("Failed to list pods per node: %v", err)
	}

	for idx := range list.Items {
		node := &list.Items[idx]
		instanceType := node.GetLabels()["node.kubernetes.io/instance-type"]
		utilisation := calculateNodeUtilisation(node, podsPerNode)
		nodes = append(nodes, NodeInfo{
			Name:              node.Name,
			Age:               time.Since(node.CreationTimestamp.Time),
			CreationTimestamp: node.CreationTimestamp.Time,
			InstanceType:      instanceType,
			Utilisation:       utilisation,
			Labels:            node.Labels,
		})
	}
	return nodes
}
