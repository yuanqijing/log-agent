package kube

import (
	"testing"
)

func TestClient_GetPodsByFieldSelector(t *testing.T) {
	client := NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})
	pods, err := client.GetPodsByFieldSelector("",
		map[string]string{
			"spec.nodeName": "sds1",
		})
	if err != nil {
		t.Error(err)
	}
	t.Log(len(pods))
}
