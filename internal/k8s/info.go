package k8s

import (
	"context"
	"time"

	"github.com/Hexta/k8s-tools/internal/containerutil"
	"github.com/Hexta/k8s-tools/internal/nodeutil"
	"github.com/Hexta/k8s-tools/internal/podutil"
	log "github.com/sirupsen/logrus"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Pods       podutil.PodInfoList
	Nodes      nodeutil.NodeInfoList
	Containers containerutil.ContainerInfoList
	ctx        context.Context
	clientset  *kubernetes.Clientset
}

func NewInfo(ctx context.Context, clientset *kubernetes.Clientset) *Info {
	return &Info{
		ctx:       ctx,
		clientset: clientset,
	}
}

func (r *Info) Fetch(nodeLabelSelector string, podLabelSelector string) {
	r.fetchPods(podLabelSelector)
	r.fetchNodes(nodeLabelSelector)
}

func (r *Info) fetchPods(labelSelector string) {
	list, err := r.clientset.CoreV1().Pods(v1.NamespaceAll).List(r.ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Failed to list containers: %v", err)
	}

	pods := make(podutil.PodInfoList, 0, len(list.Items))
	containers := make(containerutil.ContainerInfoList, 0, len(list.Items))

	for idx := range list.Items {
		pod := &list.Items[idx]

		podCPURequests := float64(0)
		podCPULimits := float64(0)

		podMemoryRequests := float64(0)
		podMemoryLimits := float64(0)

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

			containers = append(containers, &containerutil.ContainerInfo{
				Name:           podContainer.Name,
				Namespace:      pod.Namespace,
				PodName:        pod.Name,
				CPURequests:    cpuRequests,
				CPULimits:      cpuLimits,
				MemoryRequests: memoryRequests,
				MemoryLimits:   memoryLimits,
			})
		}

		pods = append(pods, &podutil.PodInfo{
			Name:              pod.Name,
			Namespace:         pod.Namespace,
			NodeName:          pod.Spec.NodeName,
			CreationTimestamp: pod.CreationTimestamp.Time,
			Labels:            pod.Labels,
			CPURequests:       podCPURequests,
			CPULimits:         podCPULimits,
			MemoryRequests:    podMemoryRequests,
			MemoryLimits:      podMemoryLimits,
		})
	}

	r.Pods = pods
	r.Containers = containers
}

func (r *Info) fetchNodes(labelSelector string) {
	log.Debugf("Listing nodes - start")
	defer log.Debugf("Listing nodes - done")

	list, err := r.clientset.CoreV1().Nodes().List(r.ctx, v1.ListOptions{
		TypeMeta:      v1.TypeMeta{},
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Panicf("Failed to list nodes: %v", err)
	}

	nodes := make(nodeutil.NodeInfoList, 0, len(list.Items))
	podsPerNode, err := listPodsPerNode(r.ctx, r.clientset)

	if err != nil {
		log.Errorf("Failed to list pods per node: %v", err)
	}

	for idx := range list.Items {
		node := &list.Items[idx]
		instanceType := node.GetLabels()["node.kubernetes.io/instance-type"]
		utilisation := calculateNodeUtilisation(node, podsPerNode)
		nodes = append(nodes, &nodeutil.NodeInfo{
			Name:              node.Name,
			Age:               time.Since(node.CreationTimestamp.Time),
			CreationTimestamp: node.CreationTimestamp.Time,
			InstanceType:      instanceType,
			Utilisation:       utilisation,
			Labels:            node.Labels,
		})
	}
	r.Nodes = nodes
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

func calculateNodeUtilisation(node *apicorev1.Node, podsPerNode map[string][]*apicorev1.Pod) nodeutil.NodeUtilisation {
	pods := podsPerNode[node.Name]

	requestsCpu := float64(0)
	requestsMemory := float64(0)

	for _, pod := range pods {
		for containerIdx := range pod.Spec.Containers {
			requestsCpu += pod.Spec.Containers[containerIdx].Resources.Requests.Cpu().AsApproximateFloat64()
			requestsMemory += pod.Spec.Containers[containerIdx].Resources.Requests.Memory().AsApproximateFloat64()
		}
	}

	return nodeutil.NodeUtilisation{
		CPU:    requestsCpu / node.Status.Allocatable.Cpu().AsApproximateFloat64(),
		Memory: requestsMemory / node.Status.Allocatable.Memory().AsApproximateFloat64(),
	}
}
