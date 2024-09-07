package pod

import (
	"context"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset, labelSelector string) InfoList {
	list, err := clientset.CoreV1().Pods(v1.NamespaceAll).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Failed to list containers: %v", err)
	}

	pods := make(InfoList, 0, len(list.Items))

	for idx := range list.Items {
		pod := &list.Items[idx]

		podCPURequests := float64(0)
		podCPULimits := float64(0)

		podMemoryRequests := float64(0)
		podMemoryLimits := float64(0)

		containers := make(container.InfoList, 0, len(pod.Spec.Containers))
		for containerIdx := range pod.Spec.Containers {
			podContainer := &pod.Spec.Containers[containerIdx]

			cpuRequests := podContainer.Resources.Requests.Cpu().AsApproximateFloat64()
			cpuLimits := podContainer.Resources.Limits.Cpu().AsApproximateFloat64()

			memoryRequests := podContainer.Resources.Requests.Memory().AsApproximateFloat64()
			memoryLimits := podContainer.Resources.Limits.Memory().AsApproximateFloat64()

			podCPURequests += cpuRequests
			podCPULimits += cpuLimits

			podMemoryRequests += memoryRequests
			podMemoryLimits += memoryLimits

			containers = append(containers, &container.Info{
				Name:           podContainer.Name,
				Namespace:      pod.Namespace,
				PodName:        pod.Name,
				CPURequests:    cpuRequests,
				CPULimits:      cpuLimits,
				MemoryRequests: memoryRequests,
				MemoryLimits:   memoryLimits,
			})
		}

		pods = append(pods, &Info{
			Name:              pod.Name,
			Namespace:         pod.Namespace,
			NodeName:          pod.Spec.NodeName,
			Containers:        containers,
			CreationTimestamp: pod.CreationTimestamp.Time,
			Labels:            pod.Labels,
			CPURequests:       podCPURequests,
			CPULimits:         podCPULimits,
			MemoryRequests:    podMemoryRequests,
			MemoryLimits:      podMemoryLimits,
		})
	}

	return pods
}
