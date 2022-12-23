package util

import (
	"github.com/go-logr/logr"
	"k8s.io/klog/v2/klogr"
)

func GetLogger() logr.Logger {
	return klogr.NewWithOptions(klogr.WithFormat(klogr.FormatKlog))
}
