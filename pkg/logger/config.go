package logger

import (
	"github.com/go-logr/logr"
	"github.com/yuanqijing/log-agent/pkg/kube"
	"github.com/yuanqijing/log-agent/pkg/util"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type Config struct {
	Logger logr.Logger

	Pods []v1.Pod
}

func (c *Config) Validate() error {
	if c.Logger.GetSink() == nil {
		return util.ErrLoggerRequired
	}
	if len(c.Pods) == 0 {
		return util.ErrPodsRequired
	}
	return nil
}

func SetupConfig() *Config {
	config := &Config{}

	// get all nodes of my node
	client := kube.APIClient
	node := kube.GetNodeName()

	pods, err := client.GetPodsByFieldSelector("", map[string]string{"spec.nodeName": node})
	if err != nil {
		panic(err)
	}

	for i, pod := range pods {
		if pod.Status.Phase != "Running" {
			continue
		}
		if strings.Contains(pod.Name, "log-agent") {
			continue
		}
		// random sample
		// TODO
		if i%10 == 0 {
			config.Pods = append(config.Pods, pod)
		}
	}
	return config
}
