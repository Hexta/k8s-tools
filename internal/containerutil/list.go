package containerutil

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListContainers(ctx context.Context, clientset *kubernetes.Clientset, labelSelector string) []ContainerInfo {
	list, err := clientset.CoreV1().Pods(v1.NamespaceAll).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		logrus.Fatalf("Failed to list containers: %v", err)
	}
	containers := make([]ContainerInfo, 0, len(list.Items))

	for idx := range list.Items {
		pod := &list.Items[idx]

		for containerIdx := range pod.Spec.Containers {
			podContainer := &pod.Spec.Containers[containerIdx]

			cpuRequests := podContainer.Resources.Requests.Cpu().AsApproximateFloat64()
			cpuLimits := podContainer.Resources.Limits.Cpu().AsApproximateFloat64()

			memoryRequests := podContainer.Resources.Requests.Memory().AsApproximateFloat64()
			memoryLimits := podContainer.Resources.Limits.Memory().AsApproximateFloat64()

			containers = append(containers, ContainerInfo{
				Name:           podContainer.Name,
				Namespace:      pod.Namespace,
				PodName:        pod.Name,
				CPURequests:    cpuRequests,
				CPULimits:      cpuLimits,
				MemoryRequests: memoryRequests,
				MemoryLimits:   memoryLimits,
			})
		}
	}

	return containers
}
