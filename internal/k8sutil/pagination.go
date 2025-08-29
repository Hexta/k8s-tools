package k8sutil

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Paginate abstracts the common Kubernetes List pagination loop using the Continue token.
// The callback receives the current ListOptions. It should perform the List call, process the items,
// and return the next continue token (list.Continue). If the next continue token is empty, iteration stops.
func Paginate(ctx context.Context, pageFn func(opts v1.ListOptions) (nextContinue string, err error)) error {
	var continueToken string
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		next, err := pageFn(v1.ListOptions{Continue: continueToken})
		if err != nil {
			return err
		}

		if continueToken = next; continueToken == "" {
			break
		}
	}

	return nil
}
