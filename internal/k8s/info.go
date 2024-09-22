package k8s

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	"github.com/Hexta/k8s-tools/internal/k8s/deployment"
	"github.com/Hexta/k8s-tools/internal/k8s/ds"
	"github.com/Hexta/k8s-tools/internal/k8s/hpa"
	k8snode "github.com/Hexta/k8s-tools/internal/k8s/node"
	k8spod "github.com/Hexta/k8s-tools/internal/k8s/pod"
	"github.com/Hexta/k8s-tools/internal/k8s/sts"
	"github.com/sethvargo/go-retry"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type Info struct {
	Pods        k8spod.InfoList
	Nodes       k8snode.InfoList
	Containers  container.InfoList
	Deployments deployment.InfoList
	DSs         ds.InfoList
	HPAs        hpa.InfoList
	STSs        sts.InfoList
	ctx         context.Context
	clientset   *kubernetes.Clientset
}

type FetchOptions struct {
	RetryInitialInterval time.Duration
	RetryJitterPercent   uint64
	RetryMaxAttempts     uint64
	RetryMaxInterval     time.Duration
}

func NewInfo(ctx context.Context, clientset *kubernetes.Clientset) *Info {
	return &Info{
		ctx:       ctx,
		clientset: clientset,
	}
}

func (r *Info) Fetch(opts FetchOptions) error {
	wg := sync.WaitGroup{}
	errorCh := make(chan error, 8)

	errorList := make([]error, 0, len(errorCh))
	go func() {
		for err := range errorCh {
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	}()

	r.startFetchFunc(r.fetchPods, &wg, opts, errorCh)
	r.startFetchFunc(r.fetchNodes, &wg, opts, errorCh)
	r.startFetchFunc(r.fetchDeployments, &wg, opts, errorCh)
	r.startFetchFunc(r.fetchHPAs, &wg, opts, errorCh)
	r.startFetchFunc(r.fetchSTSs, &wg, opts, errorCh)
	r.startFetchFunc(r.fetchDSs, &wg, opts, errorCh)

	wg.Wait()
	close(errorCh)

	return errors.Join(errorList...)
}

func (r *Info) startFetchFunc(f func(ctx context.Context) error, wg *sync.WaitGroup, opts FetchOptions, errorCh chan error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		b := newBackoff(opts)
		err := retry.Do(r.ctx, b, func(ctx context.Context) error {
			err := f(ctx)
			return retry.RetryableError(err)
		})
		errorCh <- err
	}()
}

func newBackoff(opts FetchOptions) retry.Backoff {
	b := retry.NewFibonacci(opts.RetryInitialInterval)
	b = retry.WithJitterPercent(opts.RetryJitterPercent, b)
	b = retry.WithMaxRetries(opts.RetryMaxAttempts, b)
	b = retry.WithCappedDuration(opts.RetryMaxInterval, b)

	return b
}

func (r *Info) fetchPods(ctx context.Context) error {
	log.Debugf("Listing pods - start")
	defer log.Debugf("Listing pods - done")

	var err error
	r.Pods, err = k8spod.Fetch(ctx, r.clientset)

	return err
}

func (r *Info) fetchNodes(ctx context.Context) error {
	log.Debugf("Listing nodes - start")
	defer log.Debugf("Listing nodes - done")

	var err error
	r.Nodes, err = k8snode.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchDeployments(ctx context.Context) error {
	log.Debugf("Listing deployments - start")
	defer log.Debugf("Listing deployments - done")

	var err error
	r.Deployments, err = deployment.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchHPAs(ctx context.Context) error {
	log.Debugf("Listing HPAs - start")
	defer log.Debugf("Listing HPAs - done")

	var err error
	r.HPAs, err = hpa.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchSTSs(ctx context.Context) error {
	log.Debugf("Listing STSs - start")
	defer log.Debugf("Listing STSs - done")

	var err error
	r.STSs, err = sts.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchDSs(ctx context.Context) error {
	log.Debugf("Listing DSs - start")
	defer log.Debugf("Listing DSs - done")

	var err error
	r.DSs, err = ds.Fetch(ctx, r.clientset)
	return err
}
