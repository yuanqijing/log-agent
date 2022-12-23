package kube

import (
	"fmt"
	"testing"
)

func TestLocalIP(t *testing.T) {
	ip := LocalIP()
	if ip == "" {
		t.Error("localIP() returned empty string")
	}
	fmt.Printf("localIP() returned %s\n", ip)
}
