package pod

import (
	"context"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	pods := make(InfoList, 0, 10000)

	for {
		list, err := clientset.CoreV1().Pods(v1.NamespaceAll).List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list pods: %w", err)
		}

		for idx := range list.Items {
			pod := &list.Items[idx]

			podCPURequests := float64(0)
			podCPULimits := float64(0)

			podMemoryRequests := float64(0)
			podMemoryLimits := float64(0)

			containers := getContainers(pod.Spec.Containers, pod)
			for _, c := range containers {
				podCPURequests += c.CPURequests
				podCPULimits += c.CPULimits

				podMemoryRequests += c.MemoryRequests
				podMemoryLimits += c.MemoryLimits
			}

			initContainers := getContainers(pod.Spec.InitContainers, pod)
			for _, c := range initContainers {
				podCPURequests += c.CPURequests
				podCPULimits += c.CPULimits

				podMemoryRequests += c.MemoryRequests
				podMemoryLimits += c.MemoryLimits
			}

			pods = append(pods, &Info{
				Name:        pod.Name,
				Namespace:   pod.Namespace,
				Labels:      pod.Labels,
				Annotations: pod.Annotations,

				Affinity:                      pod.Spec.Affinity,
				AutomountServiceAccountToken:  pod.Spec.AutomountServiceAccountToken,
				CPULimits:                     podCPULimits,
				CPURequests:                   podCPURequests,
				Containers:                    containers,
				CreationTimestamp:             pod.CreationTimestamp.Time,
				DNSConfig:                     pod.Spec.DNSConfig,
				DNSPolicy:                     pod.Spec.DNSPolicy,
				EnableServiceLinks:            pod.Spec.EnableServiceLinks,
				HostIPC:                       pod.Spec.HostIPC,
				HostNetwork:                   pod.Spec.HostNetwork,
				HostPID:                       pod.Spec.HostPID,
				HostUsers:                     pod.Spec.HostUsers,
				Hostname:                      pod.Spec.Hostname,
				InitContainers:                initContainers,
				IP:                            pod.Status.PodIP,
				MemoryLimits:                  podMemoryLimits,
				MemoryRequests:                podMemoryRequests,
				NodeName:                      pod.Spec.NodeName,
				NodeSelector:                  pod.Spec.NodeSelector,
				PreemptionPolicy:              pod.Spec.PreemptionPolicy,
				Priority:                      pod.Spec.Priority,
				PriorityClassName:             pod.Spec.PriorityClassName,
				RestartPolicy:                 string(pod.Spec.RestartPolicy),
				RuntimeClassName:              pod.Spec.RuntimeClassName,
				SchedulerName:                 pod.Spec.SchedulerName,
				ServiceAccountName:            pod.Spec.ServiceAccountName,
				SetHostnameAsFQDN:             pod.Spec.SetHostnameAsFQDN,
				ShareProcessNamespace:         pod.Spec.ShareProcessNamespace,
				Subdomain:                     pod.Spec.Subdomain,
				TerminationGracePeriodSeconds: pod.Spec.TerminationGracePeriodSeconds,
				Tolerations:                   getTolerations(pod.Spec.Tolerations),
			})
		}

		if continueToken = list.GetContinue(); continueToken == "" {
			break
		}
	}

	return pods, nil
}

func getTolerations(tolerations []corev1.Toleration) TolerationList {
	tolerationsList := make(TolerationList, 0, len(tolerations))
	for _, toleration := range tolerations {
		tolerationsList = append(tolerationsList, &Toleration{
			Effect:            string(toleration.Effect),
			Key:               toleration.Key,
			Operator:          string(toleration.Operator),
			TolerationSeconds: toleration.TolerationSeconds,
			Value:             toleration.Value,
		})
	}
	return tolerationsList
}

func getContainers(podContainers []corev1.Container, pod *corev1.Pod) container.InfoList {
	containers := make(container.InfoList, 0, len(podContainers))
	for containerIdx := range podContainers {
		podContainer := &podContainers[containerIdx]

		cpuRequests := podContainer.Resources.Requests.Cpu().AsApproximateFloat64()
		cpuLimits := podContainer.Resources.Limits.Cpu().AsApproximateFloat64()

		memoryRequests := podContainer.Resources.Requests.Memory().AsApproximateFloat64()
		memoryLimits := podContainer.Resources.Limits.Memory().AsApproximateFloat64()

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
	return containers
}
