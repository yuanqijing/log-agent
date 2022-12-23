package logger

import (
	"bufio"
	"context"
	"github.com/go-logr/logr"
	"github.com/yuanqijing/log-agent/apis"
	"github.com/yuanqijing/log-agent/pkg/client"
	"github.com/yuanqijing/log-agent/pkg/kube"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sync"
)

type Logger struct {
	logger logr.Logger

	// Pods defines the pods that the logger is responsible for
	// collecting logs from.
	Pods []v1.Pod
}

func NewLogger(cfg *Config) *Logger {
	err := cfg.Validate()
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: cfg.Logger,
		Pods:   cfg.Pods,
	}
}

func (l *Logger) Run(stopCh <-chan struct{}) {
	l.logger.Info("starting logger")

	// for each pod, start a goroutine to collect logs from it
	// using wait group to wait for all goroutines to finish
	// when stopCh is closed
	waitGroup := sync.WaitGroup{}
	for _, pod := range l.Pods {
		waitGroup.Add(1)
		l.logger.Info("starting logger for pod", "pod", pod.Name)
		go func(pod v1.Pod) {
			defer waitGroup.Done()
			err := l.collectLogs(pod, stopCh)
			l.logger.Error(err, "failed to collect logs", "pod", pod.Name)
		}(pod)
	}

	// wait for all goroutines to finish
	waitGroup.Wait()

	l.logger.Info("logger all finished")

	<-stopCh
}

func (l *Logger) collectLogs(pod v1.Pod, stopCh <-chan struct{}) error {
	l.logger.Info("collecting logs", "pod", pod.Name)

	ctx := context.Background()

	podLogs, err := kube.APIClient.KubeClient.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		Container:  pod.Spec.Containers[0].Name,
		Follow:     true,
		Timestamps: true,
	}).Stream(ctx)

	if err != nil {
		return err
	}

	buffer := bufio.NewReader(podLogs)
	i := 0
	for {
		line, _, err := buffer.ReadLine()
		if err != nil {
			return err
		}
		_, err = client.ESClient.Index("myindex", apis.Log{
			Source:     pod.Name,
			RawMessage: line,
		})
		if err != nil {
			l.logger.Error(err, "failed to send log to es", "pod", pod.Name)
		}
		i += 1
		if i%10 == 0 {
			klog.Infof("send %d logs to es, pod Name %s, ith message %s", i, pod.Name, line)
		}
	}
}
