package nodeutilisation

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/rodaine/table"
	log "github.com/sirupsen/logrus"
	apicorev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func PrintNodeUtilisation(ctx context.Context, kubeconfigFile string, nodepool string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	list, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{
		TypeMeta:      v1.TypeMeta{},
		LabelSelector: nodepool,
	})
	if err != nil {
		log.Fatalf("Failed to list nodes: %v", err)
	}

	nodes := make([]NodeInfo, 0, len(list.Items))

	log.Debugf("Listing Pods per node - start")
	podsPerNode, err := listPodsPerNode(ctx, clientset)
	log.Debugf("Listing Pods per node - done. Size: %v", len(podsPerNode))
	if err != nil {
		log.Errorf("Failed to list pods per node: %v", err)
	}

	for idx := range list.Items {
		node := &list.Items[idx]
		instanceType := node.GetLabels()["node.kubernetes.io/instance-type"]
		utilisation := calculateNodeResourceUtilisation(node, podsPerNode)
		nodes = append(nodes, NodeInfo{
			Name:         node.Name,
			Age:          time.Since(node.CreationTimestamp.Time),
			InstanceType: instanceType,
			Utilisation:  utilisation,
		})
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Utilisation.CPU < nodes[j].Utilisation.CPU })

	tbl := table.New("Name", "CPU", "Memory", "Type", "Age").WithWriter(log.StandardLogger().Out)
	for _, node := range nodes {
		tbl.AddRow(
			node.Name,
			fmt.Sprintf("%.3f", node.Utilisation.CPU),
			fmt.Sprintf("%.3f", node.Utilisation.Memory),
			fmt.Sprintf("%v", node.InstanceType),
			node.Age.Truncate(time.Hour),
		)
	}

	tbl.Print()
}

type NodeUtilisation struct {
	CPU    float64
	Memory float64
}

type NodeInfo struct {
	Name         string
	Age          time.Duration
	InstanceType string
	Utilisation  NodeUtilisation
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

func calculateNodeResourceUtilisation(node *apicorev1.Node, podsPerNode map[string][]*apicorev1.Pod) NodeUtilisation {
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
