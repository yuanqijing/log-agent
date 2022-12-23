package main

import (
	"fmt"
	signals "github.com/yuanqijing/log-agent/pkg/signal"
	"time"
)

func main() {
	stopChan := signals.SetupSignalHandler()

	for {
		select {
		case <-time.After(time.Second * 1):
			fmt.Println("tick")
		case <-stopChan:
			fmt.Println("stop")
			return
		}
	}
}
