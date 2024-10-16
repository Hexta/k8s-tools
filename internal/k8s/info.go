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
	Containers  container.InfoList
	DSs         ds.InfoList
	Deployments deployment.InfoList
	HPAs        hpa.InfoList
	NodeTaints  TaintList
	Nodes       k8snode.InfoList
	Pods        k8spod.InfoList
	STSs        sts.InfoList
	Taints      k8snode.TaintList
	ctx         context.Context
	clientset   *kubernetes.Clientset
}

type Taint struct {
	Effect string `db:"effect"`
	Key    string `db:"key"`
	Node   string `db:"node_name"`
	Value  string `db:"value"`
}

type TaintList []*Taint

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
	errorCh := make(chan error)

	errorList := make([]error, 0, len(errorCh))
	go func() {
		for err := range errorCh {
			if err != nil {
				errorList = append(errorList, err)
			}
		}
	}()

	r.startFetchFunc(r.fetchPods, "Pods", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchNodes, "Nodes", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchDeployments, "Deployments", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchHPAs, "HPAs", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchSTSs, "STSs", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchDSs, "DSs", &wg, opts, errorCh)

	wg.Wait()
	close(errorCh)

	return errors.Join(errorList...)
}

func (r *Info) startFetchFunc(f func(ctx context.Context) error, name string, wg *sync.WaitGroup, opts FetchOptions, errorCh chan error) {
	wg.Add(1)
	go func() {
		start := time.Now()
		defer func() {
			log.Debugf("Fetching %s - done, elapsed: %v", name, time.Since(start))
			wg.Done()
		}()
		log.Debugf("Fetching %s - start", name)
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
	var err error
	r.Pods, err = k8spod.Fetch(ctx, r.clientset)

	return err
}

func (r *Info) fetchNodes(ctx context.Context) error {
	var err error
	r.Nodes, err = k8snode.Fetch(ctx, r.clientset)
	r.NodeTaints = nodesToTaints(r.Nodes)
	return err
}

func (r *Info) fetchDeployments(ctx context.Context) error {
	var err error
	r.Deployments, err = deployment.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchHPAs(ctx context.Context) error {
	var err error
	r.HPAs, err = hpa.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchSTSs(ctx context.Context) error {
	var err error
	r.STSs, err = sts.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchDSs(ctx context.Context) error {
	var err error
	r.DSs, err = ds.Fetch(ctx, r.clientset)
	return err
}

func nodesToTaints(nodes k8snode.InfoList) TaintList {
	taints := make(TaintList, 0, len(nodes))
	for _, node := range nodes {
		for _, taint := range node.Taints {
			taints = append(taints, &Taint{
				Node:   node.Name,
				Key:    taint.Key,
				Value:  taint.Value,
				Effect: taint.Effect,
			})
		}
	}

	return taints
}
