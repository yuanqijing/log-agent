package util

import "fmt"

var (
	ErrElectorConfigRequired      = fmt.Errorf("elector config is required")
	ErrLoggerRequired             = fmt.Errorf("logger is required")
	ErrLeaseLockNameRequired      = fmt.Errorf("lease lock name is required")
	ErrLeaseLockNamespaceRequired = fmt.Errorf("lease lock namespace is required")
	ErrPodsRequired               = fmt.Errorf("pod id is required")
)
