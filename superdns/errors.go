package superdns

import (
	"errors"
)

// errors
var (
	ErrClusterNotFound    = errors.New("can not find cluster")
	ErrNoAvalibaleCluster = errors.New("no available cluster")
)
