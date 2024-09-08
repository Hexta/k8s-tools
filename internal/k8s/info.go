package k8s

import (
	"context"
	"errors"
	"sync"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	"github.com/Hexta/k8s-tools/internal/k8s/deployment"
	"github.com/Hexta/k8s-tools/internal/k8s/hpa"
	k8snode "github.com/Hexta/k8s-tools/internal/k8s/node"
	k8spod "github.com/Hexta/k8s-tools/internal/k8s/pod"
	"github.com/Hexta/k8s-tools/internal/k8s/sts"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Pods        k8spod.InfoList
	Nodes       k8snode.InfoList
	Containers  container.InfoList
	Deployments deployment.InfoList
	HPAs        hpa.InfoList
	STSs        sts.InfoList
	ctx         context.Context
	clientset   *kubernetes.Clientset
}

func NewInfo(ctx context.Context, clientset *kubernetes.Clientset) *Info {
	return &Info{
		ctx:       ctx,
		clientset: clientset,
	}
}

func (r *Info) Fetch(nodeLabelSelector string, podLabelSelector string) error {
	wg := sync.WaitGroup{}

	errorCh := make(chan error, 8)

	wg.Add(1)
	go func() { defer wg.Done(); err := r.fetchPods(); errorCh <- err }()

	wg.Add(1)
	go func() { defer wg.Done(); err := r.fetchNodes(); errorCh <- err }()

	wg.Add(1)
	go func() { defer wg.Done(); err := r.fetchDeployments(); errorCh <- err }()

	wg.Add(1)
	go func() { defer wg.Done(); err := r.fetchHPAs(); errorCh <- err }()

	wg.Add(1)
	go func() { defer wg.Done(); err := r.fetchSTSs(); errorCh <- err }()

	errorList := make([]error, 0, len(errorCh))
	go func() {
		for err := range errorCh {
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	}()

	wg.Wait()
	close(errorCh)

	return errors.Join(errorList...)
}

func (r *Info) fetchPods() error {
	log.Debugf("Listing pods - start")
	defer log.Debugf("Listing pods - done")

	var err error
	r.Pods, err = k8spod.Fetch(r.ctx, r.clientset)

	return err
}

func (r *Info) fetchNodes() error {
	log.Debugf("Listing nodes - start")
	defer log.Debugf("Listing nodes - done")

	var err error
	r.Nodes, err = k8snode.Fetch(r.ctx, r.clientset)
	return err
}

func (r *Info) fetchDeployments() error {
	log.Debugf("Listing deployments - start")
	defer log.Debugf("Listing deployments - done")

	var err error
	r.Deployments, err = deployment.Fetch(r.ctx, r.clientset)
	return err
}

func (r *Info) fetchHPAs() error {
	log.Debugf("Listing HPAs - start")
	defer log.Debugf("Listing HPAs - done")

	var err error
	r.HPAs, err = hpa.Fetch(r.ctx, r.clientset)
	return err
}

func (r *Info) fetchSTSs() error {
	log.Debugf("Listing STSs - start")
	defer log.Debugf("Listing STSs - done")

	var err error
	r.STSs, err = sts.Fetch(r.ctx, r.clientset)
	return err
}
