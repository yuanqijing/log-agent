package kube

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

func (c *Client) PatchPodLabels(namespace, name string, label map[string]string) error {
	_, err := c.KubeClient.CoreV1().Pods(namespace).Patch(
		context.Background(),
		name,
		types.MergePatchType,
		[]byte(`{"metadata":{"labels":`+Map2Json(label)+`}}`),
		metav1.PatchOptions{},
	)
	if err != nil {
		klog.Errorf("PatchPodLabels error: %v", err)
		return err
	}
	return nil
}

func (c *Client) GetPod(namespace, name string) (*v1.Pod, error) {
	pod, err := c.KubeClient.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("GetPod error: %v", err)
		return nil, err
	}
	return pod, nil
}

func (c *Client) GetPodLabels(namespace, name string) (map[string]string, error) {
	pod, err := c.KubeClient.CoreV1().Pods(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("GetPodLabels error: %v", err)
		return nil, err
	}
	return pod.Labels, nil
}

func (c *Client) GetPodsBySelector(namespace string, selector map[string]string) ([]v1.Pod, error) {
	pods, err := c.KubeClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: Map2Str(selector),
	})
	if err != nil {
		klog.Errorf("GetPodsBySelector error: %v", err)
		return nil, err
	}
	return pods.Items, nil
}

func (c *Client) GetPodsByFieldSelector(namespace string, fieldSelector map[string]string) ([]v1.Pod, error) {
	pods, err := c.KubeClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		FieldSelector: Map2Str(fieldSelector),
	})
	if err != nil {
		klog.Errorf("GetPodsByFieldSelector error: %v", err)
		return nil, err
	}
	return pods.Items, nil
}

func (c *Client) GetPodsIPsBySelector(namespace string, selector map[string]string) ([]string, error) {
	pods, err := c.KubeClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: Map2Str(selector),
	})
	if err != nil {
		klog.Errorf("GetPodsIPsBySelector error: %v", err)
		return nil, err
	}
	var ips []string
	for _, pod := range pods.Items {
		ips = append(ips, pod.Status.PodIP)
	}
	return ips, nil
}
