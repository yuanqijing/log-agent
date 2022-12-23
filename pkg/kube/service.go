package kube

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

// service.go implements methods for modify a kubernetes service.

// PatchServiceLabels by default uses strategic merge patch.
// if label k v does ont exist, it will add. if k v exist, it will update.
func (c *Client) PatchServiceLabels(namespace, name string, labels map[string]string) error {
	_, err := c.KubeClient.CoreV1().Services(namespace).Patch(
		context.Background(),
		name,
		types.MergePatchType,
		[]byte(`{"metadata":{"labels":`+Map2Json(labels)+`}}`),
		metav1.PatchOptions{},
	)
	if err != nil {
		klog.Errorf("PatchServiceLabels error: %v", err)
		return err
	}
	return nil
}

func (c *Client) CreateIfNotExistService(namespace string, service *v1.Service) error {
	// check if service exist
	_, err := c.KubeClient.CoreV1().Services(namespace).Get(context.Background(), service.Name, metav1.GetOptions{})
	if err == nil {
		return nil
	}
	_, err = c.KubeClient.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("CreateIfNotExistService error: %v", err)
		return err
	}
	return nil
}

func (c *Client) GetService(namespace, name string) (*v1.Service, error) {
	// TODO: It is good to use pointer ?
	service, err := c.KubeClient.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("GetService error: %v", err)
		return nil, err
	}
	return service, nil
}

// GetServiceLabels returns labels of a kubernetes service.
func (c *Client) GetServiceLabels(namespace, name string) (map[string]string, error) {
	service, err := c.KubeClient.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("GetServiceLabels error: %v", err)
		return nil, err
	}
	return service.Labels, nil
}

// PatchServiceSelector by default uses MergePatchType patch.
func (c *Client) PatchServiceSelector(namespace, name string, selector map[string]string) error {
	_, err := c.KubeClient.CoreV1().Services(namespace).Patch(
		context.Background(),
		name,
		types.MergePatchType,
		[]byte(`{"spec":{"selector":`+Map2Json(selector)+`}}`),
		metav1.PatchOptions{},
	)
	if err != nil {
		klog.Errorf("PatchServiceSelector error: %v", err)
		return err
	}
	return nil
}

// GetServiceSelector returns selector of a kubernetes service.
func (c *Client) GetServiceSelector(namespace, name string) (map[string]string, error) {
	service, err := c.KubeClient.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("GetServiceSelector error: %v", err)
		return nil, err
	}
	return service.Spec.Selector, nil
}

// SimplifyService returns a simplified service.
func SimplifyService(service *v1.Service) *v1.Service {
	service.ObjectMeta.UID = ""
	service.ObjectMeta.ResourceVersion = ""
	service.ObjectMeta.CreationTimestamp = metav1.Time{}
	service.ObjectMeta.DeletionTimestamp = nil
	service.ObjectMeta.DeletionGracePeriodSeconds = nil
	service.Spec.ClusterIP = ""
	service.Spec.ClusterIPs = nil
	service.Spec.ExternalIPs = nil
	service.Spec.LoadBalancerSourceRanges = nil
	service.Spec.ExternalName = ""
	service.Spec.ExternalTrafficPolicy = ""
	service.Spec.HealthCheckNodePort = 0
	service.Spec.PublishNotReadyAddresses = false
	service.Spec.SessionAffinity = ""
	service.Spec.SessionAffinityConfig = nil
	service.Spec.ExternalTrafficPolicy = ""
	return &v1.Service{
		ObjectMeta: service.ObjectMeta,
		Spec:       service.Spec,
	}
}
