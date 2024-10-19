package k8s

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	"github.com/Hexta/k8s-tools/internal/k8s/deployment"
	"github.com/Hexta/k8s-tools/internal/k8s/ds"
	"github.com/Hexta/k8s-tools/internal/k8s/fetch"
	"github.com/Hexta/k8s-tools/internal/k8s/hpa"
	k8snode "github.com/Hexta/k8s-tools/internal/k8s/node"
	k8spod "github.com/Hexta/k8s-tools/internal/k8s/pod"
	k8sservice "github.com/Hexta/k8s-tools/internal/k8s/service"
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
	Nodes       k8snode.InfoList
	Pods        k8spod.InfoList
	Services    k8sservice.InfoList
	STSs        sts.InfoList
	Taints      TaintList
	Tolerations TolerationList
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

type Toleration struct {
	Effect            string `db:"effect"`
	Key               string `db:"key"`
	Pod               string `db:"pod_name"`
	Value             string `db:"value"`
	TolerationSeconds *int64 `db:"toleration_seconds"`
	Operator          string `db:"operator"`
}

type TolerationList []*Toleration

func NewInfo(ctx context.Context, clientset *kubernetes.Clientset) *Info {
	return &Info{
		ctx:       ctx,
		clientset: clientset,
	}
}

func (r *Info) Fetch(opts fetch.Options) error {
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

	start := time.Now()
	log.Debug("Fetching K8s info")

	r.startFetchFunc(r.fetchPods, "Pods", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchNodes, "Nodes", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchDeployments, "Deployments", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchHPAs, "HPAs", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchSTSs, "STSs", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchDSs, "DSs", &wg, opts, errorCh)
	r.startFetchFunc(r.fetchServices, "Services", &wg, opts, errorCh)

	wg.Wait()
	close(errorCh)

	log.Debug("Fetching K8s info: done, elapsed: ", time.Since(start))

	return errors.Join(errorList...)
}

func (r *Info) startFetchFunc(f func(ctx context.Context, opts fetch.Options) error, name string, wg *sync.WaitGroup, opts fetch.Options, errorCh chan error) {
	wg.Add(1)
	go func() {
		start := time.Now()
		defer func() {
			log.Debugf("Fetching %s: done, elapsed: %v", name, time.Since(start))
			wg.Done()
		}()
		log.Debugf("Fetching %s", name)
		b := newBackoff(opts)
		err := retry.Do(r.ctx, b, func(ctx context.Context) error {
			err := f(ctx, opts)
			return retry.RetryableError(err)
		})
		errorCh <- err
	}()
}

func newBackoff(opts fetch.Options) retry.Backoff {
	b := retry.NewFibonacci(opts.RetryInitialInterval)
	b = retry.WithJitterPercent(opts.RetryJitterPercent, b)
	b = retry.WithMaxRetries(opts.RetryMaxAttempts, b)
	b = retry.WithCappedDuration(opts.RetryMaxInterval, b)

	return b
}

func (r *Info) fetchPods(ctx context.Context, _ fetch.Options) error {
	var err error
	r.Pods, err = k8spod.Fetch(ctx, r.clientset)
	r.Tolerations = podsToTolerations(r.Pods)

	return err
}

func (r *Info) fetchNodes(ctx context.Context, opts fetch.Options) error {
	var err error
	r.Nodes, err = k8snode.Fetch(ctx, r.clientset, opts)
	r.Taints = nodesToTaints(r.Nodes)
	return err
}

func (r *Info) fetchDeployments(ctx context.Context, _ fetch.Options) error {
	var err error
	r.Deployments, err = deployment.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchHPAs(ctx context.Context, _ fetch.Options) error {
	var err error
	r.HPAs, err = hpa.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchSTSs(ctx context.Context, _ fetch.Options) error {
	var err error
	r.STSs, err = sts.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchDSs(ctx context.Context, _ fetch.Options) error {
	var err error
	r.DSs, err = ds.Fetch(ctx, r.clientset)
	return err
}

func (r *Info) fetchServices(ctx context.Context, _ fetch.Options) error {
	var err error
	r.Services, err = k8sservice.Fetch(ctx, r.clientset)
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

func podsToTolerations(pods k8spod.InfoList) TolerationList {
	tolerations := make(TolerationList, 0, len(pods))
	for _, pod := range pods {
		for _, toleration := range pod.Tolerations {
			tolerations = append(tolerations, &Toleration{
				Pod:               pod.Name,
				Key:               toleration.Key,
				Value:             toleration.Value,
				TolerationSeconds: toleration.TolerationSeconds,
				Operator:          toleration.Operator,
			})
		}
	}

	return tolerations
}
