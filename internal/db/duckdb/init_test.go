package duckdb_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hexta/k8s-tools/internal/db"
	"github.com/Hexta/k8s-tools/internal/k8s"
	"github.com/Hexta/k8s-tools/internal/k8s/fetch"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	fakeapiextensions "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestDBInit(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resources := []runtime.Object{
		&v1.Node{
			ObjectMeta: metav1.ObjectMeta{Name: "node1"},
		},
	}

	// Initialize fake clients
	clientSet := k8sfake.NewSimpleClientset(resources...)
	dynamicClient := dynamicfake.NewSimpleDynamicClient(scheme.Scheme)
	apiExtClient := fakeapiextensions.NewSimpleClientset()

	k8sInfo := k8s.NewInfo(ctx, clientSet, dynamicClient, apiExtClient)

	// Create a temporary cache directory for the database
	dbDir := t.TempDir()

	// Fetch info (using fake clients)
	err := k8sInfo.Fetch(fetch.Options{RetryInitialInterval: 100 * time.Millisecond, RetryMaxAttempts: 0})
	require.NoErrorf(t, err, "Failed to fetch K8s info")

	// Initialize database
	err = db.InitDB(ctx, dbDir, k8sInfo)
	require.NoErrorf(t, err, "Failed to init DB")

	// Query the database
	data, err := db.Query(ctx, dbDir, "SELECT * FROM k8s.nodes")
	require.NoErrorf(t, err, "Failed to query DB")
	require.NotNil(t, data)
}
