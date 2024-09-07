package k8s

import (
	"context"
	"sync"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	k8snode "github.com/Hexta/k8s-tools/internal/k8s/node"
	k8spod "github.com/Hexta/k8s-tools/internal/k8s/pod"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Pods       k8spod.InfoList
	Nodes      k8snode.InfoList
	Containers container.InfoList
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
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() { r.fetchPods(podLabelSelector); wg.Done() }()

	wg.Add(1)
	go func() { r.fetchNodes(nodeLabelSelector); wg.Done() }()

	wg.Wait()
}

func (r *Info) fetchPods(labelSelector string) {
	log.Debugf("Listing pods - start")
	defer log.Debugf("Listing pods - done")

	r.Pods = k8spod.Fetch(r.ctx, r.clientset, labelSelector)
}

func (r *Info) fetchNodes(labelSelector string) {
	log.Debugf("Listing nodes - start")
	defer log.Debugf("Listing nodes - done")

	r.Nodes = k8snode.Fetch(r.ctx, r.clientset, labelSelector)
}
