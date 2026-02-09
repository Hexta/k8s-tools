package k8s

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(restConfig *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(restConfig)
}

func GetRestConfig(kubeconfigFile string, kubeContext string) (*rest.Config, error) {
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

	return clientConfig.ClientConfig()
}

func GetDynamicClient(restConfig *rest.Config) (*dynamic.DynamicClient, error) {
	return dynamic.NewForConfig(restConfig)
}

func GetAPIExtensionsClient(restConfig *rest.Config) (apiextensionsclient.Interface, error) {
	return apiextensionsclient.NewForConfig(restConfig)
}
