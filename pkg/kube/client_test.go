package kube

import "testing"

func TestNewClient(t *testing.T) {
	_ = NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})
}

func TestClient_PatchServiceLabels(t *testing.T) {
	client := NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})

	lables := map[string]string{
		"app":                         "redis-t",
		"component":                   "redis-t",
		"harmonycloud.cn/statefulset": "redis-t",
		"middleware":                  "redis",
		"nephele/user":                "admin",
	}

	if err := client.PatchServiceLabels("qijing", "redis-t-readonly", lables); err != nil {
		t.Error(err)
	}
}

func TestClient_PatchServiceSelector(t *testing.T) {
	client := NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})

	selector := map[string]string{
		"app": "redis-t",
	}

	if err := client.PatchServiceSelector("qijing", "redis-t-readonly", selector); err != nil {
		t.Error(err)
	}
}

func TestClient_PatchPodLabels(t *testing.T) {
	client := NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})

	lables := map[string]string{
		"app.kubernetes.io/instance":              "gossip-agent-2022-08-10-b85c3f9",
		"app.kubernetes.io/name":                  "agent",
		"gossip.middleware.com/gossip-agent-id":   "gossip-agent-2022-08-10-b85c3f9-5b59b9b587-bnkdl",
		"gossip.middleware.com/gossip-agent-type": "daemon",
		"pod-template-hash":                       "5b59b9b587",
	}

	if err := client.PatchPodLabels("qijing", "gossip-agent-2022-08-10-b85c3f9-5b59b9b587-fnw5q", lables); err != nil {
		t.Error(err)
	}
}

func TestClient_GetPodsBySelector(t *testing.T) {
	client := NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})

	pods, err := client.GetPodsBySelector("qijing", map[string]string{
		"app": "redis-t",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(pods)
}

func TestClient_GetPodsIPsBySelector(t *testing.T) {
	client := NewClient(Config{
		KubeConfig: "/Users/mac/.kube/config",
	})

	pods, err := client.GetPodsIPsBySelector("qijing", map[string]string{
		"app": "redis-t",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(pods)
}

func TestMap2Json(t *testing.T) {
	m := map[string]string{
		"app": "redis-t",
	}
	t.Log(Map2Json(m))
}

func TestMap2Str(t *testing.T) {
	m := map[string]string{
		"app": "redis-t",
		"env": "dev",
	}
	t.Log(Map2Str(m))
}
