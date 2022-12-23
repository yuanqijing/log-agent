package kube

import (
	"fmt"
)

type Config struct {
	KubeConfig string `json:"kubeconfig"`
}

// String is used to print the config, for debugging.
func (c *Config) String() string {
	return fmt.Sprintf(""+
		"Config: {\n"+
		"  KubeConfig: %s\n"+
		"} \n", c.KubeConfig)
}
