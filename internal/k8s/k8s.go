package k8s

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(kubeconfigFile string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
	if err != nil {
		log.Panicf("Failed to create config: %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}

	return clientSet
}
