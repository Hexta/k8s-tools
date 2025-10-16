package k8s

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(kubeconfigFile string, kubeContext string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.LoadFromFile(kubeconfigFile)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	if kubeContext != "" {
		if _, exists := config.Contexts[kubeContext]; !exists {
			return nil, fmt.Errorf("context %s does not exist", kubeContext)
		}
	}

	clientConfig := clientcmd.NewNonInteractiveClientConfig(*config, kubeContext, &clientcmd.ConfigOverrides{}, nil)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get client config: %w", err)
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return clientSet, nil
}
