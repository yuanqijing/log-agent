package main

import (
	daemon2 "github.com/yuanqijing/log-agent/pkg/daemon"
	signals "github.com/yuanqijing/log-agent/pkg/signal"
)

func main() {

	stopCh := signals.SetupSignalHandler()

	// ...
	daemon := daemon2.NewDaemon()
	daemon.Run(stopCh)
}
