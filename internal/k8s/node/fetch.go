package node

import (
	"context"
	"fmt"
	"time"

	"github.com/Hexta/k8s-tools/internal/k8s/fetch"
	"github.com/Hexta/k8s-tools/internal/k8sutil"
	log "github.com/sirupsen/logrus"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const instanceTypeLabel = "node.kubernetes.io/instance-type"

func Fetch(ctx context.Context, clientset *kubernetes.Clientset, opts fetch.Options) (InfoList, error) {
	nodes := make(InfoList, 0, 10000)

	podsPerNode, err := listPodsPerNode(ctx, clientset)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods per node: %w", err)
	}

	err = k8sutil.Paginate(ctx, func(lo v1.ListOptions) (string, error) {
		// Preserve label selector while paginating
		lo.LabelSelector = opts.LabelSelector
		list, err := clientset.CoreV1().Nodes().List(ctx, lo)
		if err != nil {
			return "", fmt.Errorf("failed to list nodes: %w", err)
		}

		for idx := range list.Items {
			node := &list.Items[idx]
			instanceType := node.GetLabels()[instanceTypeLabel]
			utilisation := calculateNodeUtilisation(node, podsPerNode)

			addrMap := getAddressMap(node.Status.Addresses)
			taints := getTaints(node.Spec.Taints)
			images := getImages(node)

			nodes = append(nodes, &Info{
				Address:                 addrMap,
				Age:                     time.Since(node.CreationTimestamp.Time),
				AllocatableCPU:          node.Status.Allocatable.Cpu().AsApproximateFloat64(),
				AllocatableMemory:       node.Status.Allocatable.Memory().AsApproximateFloat64(),
				Annotations:             node.Annotations,
				Architecture:            node.Status.NodeInfo.Architecture,
				CapacityCPU:             node.Status.Capacity.Cpu().AsApproximateFloat64(),
				CapacityMemory:          node.Status.Capacity.Memory().AsApproximateFloat64(),
				ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
				CPUUtilisation:          utilisation.CPU,
				CreationTimestamp:       node.CreationTimestamp.Time,
				Images:                  images,
				InstanceType:            instanceType,
				KernelVersion:           node.Status.NodeInfo.KernelVersion,
				KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
				Labels:                  node.Labels,
				MemoryUtilisation:       utilisation.Memory,
				Name:                    node.Name,
				OperatingSystem:         node.Status.NodeInfo.OperatingSystem,
				OSImage:                 node.Status.NodeInfo.OSImage,
				Taints:                  taints,
			})
		}

		return list.Continue, nil
	})

	return nodes, err
}

func getTaints(taints []apicorev1.Taint) TaintList {
	taintList := make([]*Taint, 0, len(taints))

	for _, taint := range taints {
		taintList = append(taintList, &Taint{
			Effect: string(taint.Effect),
			Key:    taint.Key,
			Value:  taint.Value,
		})
	}

	return taintList
}

func getImages(node *apicorev1.Node) ImageList {
	imageList := make(ImageList, 0, len(node.Status.Images))

	for _, image := range node.Status.Images {
		imageList = append(imageList, &Image{
			Names: image.Names,
			Node:  node.Name,
			Size:  image.Size(),
		})
	}

	return imageList
}

func getAddressMap(addresses []apicorev1.NodeAddress) map[string]string {
	addrMap := make(map[string]string, 4)

	for _, address := range addresses {
		addrMap[string(address.Type)] = address.Address
	}

	return addrMap
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

func calculateNodeUtilisation(node *apicorev1.Node, podsPerNode map[string][]*apicorev1.Pod) Utilisation {
	pods := podsPerNode[node.Name]

	requestsCpu := float64(0)
	requestsMemory := float64(0)

	for _, pod := range pods {
		for containerIdx := range pod.Spec.Containers {
			requestsCpu += pod.Spec.Containers[containerIdx].Resources.Requests.Cpu().AsApproximateFloat64()
			requestsMemory += pod.Spec.Containers[containerIdx].Resources.Requests.Memory().AsApproximateFloat64()
		}
	}

	return Utilisation{
		CPU:    requestsCpu / node.Status.Allocatable.Cpu().AsApproximateFloat64(),
		Memory: requestsMemory / node.Status.Allocatable.Memory().AsApproximateFloat64(),
	}
}
