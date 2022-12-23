package elector

import (
	"github.com/go-logr/logr"
	"github.com/yuanqijing/log-agent/pkg/util"
)

type Config struct {
	Logger logr.Logger

	// leaseLockName is the lease lock resource name
	LeaseLockName string `json:"leaseLockName,omitempty"`

	// leaseLockNamespace is the lease lock resource namespace
	LeaseLockNamespace string `json:"leaseLockNamespace,omitempty"`

	// kubeconfig path to kubeconfig file
	KubeConfig string `json:"kubeconfig,omitempty"`
}

func (c *Config) Validate() error {
	// check if logger is empty
	if c.Logger.GetSink() == nil {
		return util.ErrLoggerRequired
	}

	if c.LeaseLockName == "" {
		return util.ErrLeaseLockNameRequired
	}
	if c.LeaseLockNamespace == "" {
		return util.ErrLeaseLockNamespaceRequired
	}
	return nil
}
