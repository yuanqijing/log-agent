package daemon

import (
	"github.com/yuanqijing/log-agent/pkg/elector"
	"github.com/yuanqijing/log-agent/pkg/logger"
	"k8s.io/klog/v2"
	"time"
)

type Daemon struct {
	// SampleRate is the rate at which the daemon will sample the cluster
	SampleRate int

	// Elector is the elector
	Elector *elector.Elector

	// Logger is the logger
	Logger *logger.Logger
}

func NewDaemon() *Daemon {
	// setup configuration
	config, err := SetupConfig()
	if err != nil {
		panic(err)
	}

	// setup elector
	e := elector.NewElector(config.ElectorConfig)

	// setup logger
	l := logger.NewLogger(config.LoggerConfig)

	return &Daemon{
		Elector: e,
		Logger:  l,
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
			d.runAsLeader(stopCh)
		} else {
			klog.Infof("I am not the leader")
		}
		<-time.After(5 * time.Second)
	}
}

func (d *Daemon) runAsLeader(stopCh <-chan struct{}) {
	klog.Infof("I am the leader")
	go d.Logger.Run(stopCh)
	<-stopCh
	klog.Infof("Leader daemon is shutting down...")
}
