package nodeutil

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func ListNodes(ctx context.Context, kubeconfigFile string, labelSelector string) []NodeInfo {
	log.Debugf("Listing nodes - start")
	defer log.Debugf("Listing nodes - done")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		log.Panicf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}

	list, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{
		TypeMeta:      v1.TypeMeta{},
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Panicf("Failed to list nodes: %v", err)
	}

	nodes := make([]NodeInfo, 0, len(list.Items))
	podsPerNode, err := listPodsPerNode(ctx, clientset)

	if err != nil {
		log.Errorf("Failed to list pods per node: %v", err)
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

func listPodsPerNode(ctx context.Context, clientset *kubernetes.Clientset) (map[string][]*apicorev1.Pod, error) {
	log.Debugf("Listing Pods per node - start")
	defer log.Debugf("Listing Pods per node - done")

	podsPerNode := make(map[string][]*apicorev1.Pod)

	list, err := clientset.CoreV1().Pods("").List(ctx, v1.ListOptions{})
	if err != nil {
		log.Errorf("Failed to list pods: %v", err)
		return podsPerNode, err
	}

	for i := range list.Items {
		if _, ok := podsPerNode[list.Items[i].Spec.NodeName]; !ok {
			podsPerNode[list.Items[i].Spec.NodeName] = make([]*apicorev1.Pod, 0, 8)
		}
		podsPerNode[list.Items[i].Spec.NodeName] = append(podsPerNode[list.Items[i].Spec.NodeName], &list.Items[i])
	}

	return podsPerNode, nil
}
