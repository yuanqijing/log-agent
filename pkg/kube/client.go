package kube

import (
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var APIClient = NewClient(Config{})

type Client struct {
	KubeConfig string

	KubeClient *clientset.Clientset
}

func NewClient(cfg Config) *Client {
	restConfig, err := buildConfig(cfg.KubeConfig)
	if err != nil {
		klog.Fatal(err)
	}
	// if err, then die
	kubeClient := clientset.NewForConfigOrDie(restConfig)

	return &Client{
		KubeConfig: cfg.KubeConfig,
		KubeClient: kubeClient,
	}
}

// NewClientWithClients
// TODO: consider using a clientset if already have a clientset instance
func NewClientWithClients(client *clientset.Clientset) *Client {
	return &Client{
		KubeClient: client,
	}
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
