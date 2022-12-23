package util

import "testing"

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	if logger.GetSink() == nil {
		t.Error("logger sink is nil")
	}
	anotherLogger := logger.WithName("another")
	if anotherLogger.GetSink() == nil {
		t.Error("another logger sink is nil")
	}
}
