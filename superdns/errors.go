package superdns

import (
	"errors"
)

// errors
var (
	ErrClusterNotFound     = errors.New("can not find cluster")
	ErrInvalidDestinations = errors.New("invalid destinations")
)
