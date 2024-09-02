package podutil

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListPods(ctx context.Context, clientset *kubernetes.Clientset, labelSelector string) []PodInfo {
	list, err := clientset.CoreV1().Pods(v1.NamespaceAll).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		logrus.Fatalf("Failed to list pods: %v", err)
	}
	pods := make([]PodInfo, 0, len(list.Items))

	for idx := range list.Items {
		pod := &list.Items[idx]

		cpuRequests := float64(0)
		memoryRequests := float64(0)

		for containerIdx := range pod.Spec.Containers {
			cpuRequests += pod.Spec.Containers[containerIdx].Resources.Requests.Cpu().AsApproximateFloat64()
			memoryRequests += pod.Spec.Containers[containerIdx].Resources.Requests.Memory().AsApproximateFloat64()
		}

		pods = append(pods, PodInfo{
			Name:              pod.Name,
			Namespace:         pod.Namespace,
			NodeName:          pod.Spec.NodeName,
			CreationTimestamp: pod.CreationTimestamp.Time,
			Labels:            pod.Labels,
			CPURequests:       cpuRequests,
			MemoryRequests:    memoryRequests,
		})
	}

	return pods
}
