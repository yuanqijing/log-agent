package daemon

import (
	"github.com/yuanqijing/log-agent/pkg/elector"
	"k8s.io/klog/v2"
	"time"
)

type Daemon struct {
	// SampleRate is the rate at which the daemon will sample the cluster
	SampleRate int

	// Elector is the elector
	Elector *elector.Elector
}

func NewDaemon() *Daemon {
	// setup configuration
	config, err := SetupConfig()
	if err != nil {
		panic(err)
	}

	// setup elector
	e := elector.NewElector(config.ElectorConfig)

	return &Daemon{
		Elector: e,
	}
}

func (d *Daemon) Run(stopCh <-chan struct{}) {
	// elector routine
	go d.Elector.Run(stopCh)

	// run main routine
	go d.run(stopCh)

	<-stopCh
	klog.Infof("Daemon is shutting down...")
}

func (d *Daemon) run(stopCh <-chan struct{}) {
	klog.Infof("Daemon is running...")
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		isLeader := d.Elector.GetLeader()
		if isLeader {
			klog.Infof("I am the leader")
			// TODO: do leader stuff
		} else {
			klog.Infof("I am not the leader")
		}

		<-time.After(5 * time.Second)
	}
}
