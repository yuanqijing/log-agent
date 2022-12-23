package daemon

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/yuanqijing/log-agent/pkg/elector"
	"github.com/yuanqijing/log-agent/pkg/util"
	"io/ioutil"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

type Config struct {
	ElectorConfig *elector.Config `json:"electorConfig,omitempty"`
}

func SetupConfig() (*Config, error) {
	config := &Config{}

	// --config flag is set
	path := "/etc/log-agent/config/config.yaml"

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(file, config); err != nil {
		panic(err)
	}

	klog.Infof("config: %s", spew.Sdump(config))

	baseLogger := util.GetLogger()
	config.ElectorConfig.Logger = baseLogger.WithName("elector")

	if err = config.Validate(); err != nil {
		panic(err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.ElectorConfig == nil {
		return util.ErrElectorConfigRequired
	}
	if err := c.ElectorConfig.Validate(); err != nil {
		return err
	}
	return nil
}
