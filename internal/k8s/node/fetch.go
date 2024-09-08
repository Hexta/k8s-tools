package node

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset, labelSelector string) (InfoList, error) {
	var continueToken string
	nodes := make(InfoList, 0, 10000)

	for {
		log.Debugf("Listing nodes with label selector %q, continue token %q", labelSelector, continueToken)
		list, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{
			TypeMeta:      v1.TypeMeta{},
			LabelSelector: labelSelector,
			Continue:      continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list nodes: %v", err)
		}

		podsPerNode, err := listPodsPerNode(ctx, clientset)

		if err != nil {
			return nil, fmt.Errorf("failed to list pods per node: %v", err)
		}

		for idx := range list.Items {
			node := &list.Items[idx]
			instanceType := node.GetLabels()["node.kubernetes.io/instance-type"]
			utilisation := calculateNodeUtilisation(node, podsPerNode)
			nodes = append(nodes, &Info{
				Name:              node.Name,
				Age:               time.Since(node.CreationTimestamp.Time),
				CreationTimestamp: node.CreationTimestamp.Time,
				InstanceType:      instanceType,
				Utilisation:       utilisation,
				Labels:            node.Labels,
			})
		}

		if continueToken = list.GetContinue(); continueToken == "" {
			break
		}
	}

	return nodes, nil
}

func listPodsPerNode(ctx context.Context, clientset *kubernetes.Clientset) (map[string][]*apicorev1.Pod, error) {
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

func calculateNodeUtilisation(node *apicorev1.Node, podsPerNode map[string][]*apicorev1.Pod) NodeUtilisation {
	pods := podsPerNode[node.Name]

	requestsCpu := float64(0)
	requestsMemory := float64(0)

	for _, pod := range pods {
		for containerIdx := range pod.Spec.Containers {
			requestsCpu += pod.Spec.Containers[containerIdx].Resources.Requests.Cpu().AsApproximateFloat64()
			requestsMemory += pod.Spec.Containers[containerIdx].Resources.Requests.Memory().AsApproximateFloat64()
		}
	}

	return NodeUtilisation{
		CPU:    requestsCpu / node.Status.Allocatable.Cpu().AsApproximateFloat64(),
		Memory: requestsMemory / node.Status.Allocatable.Memory().AsApproximateFloat64(),
	}
}
