package common

import (
	"github.com/aerospike/aerospike-client-go/types"
	"github.com/pkg/errors"
	"strings"
)

//NodeDownError node down error
var NodeDownError = errors.New("node is down")

//IsKeyNotFound returns true if key not found error
func IsKeyNotFound(err error) bool {
	if err == nil {
		return false
	}
	aeroError, ok := err.(types.AerospikeError)
	if !ok {
		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
		if aeroError, ok = err.(types.AerospikeError); !ok {
			return false
		}

	}
	return aeroError.ResultCode() == types.KEY_NOT_FOUND_ERROR
}

//IsTimeout returns true if timeout error
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	aeroError, ok := err.(types.AerospikeError)
	if !ok {
		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
		if aeroError, ok = err.(types.AerospikeError); !ok {
			return false
		}

	}
	return aeroError.ResultCode() == types.TIMEOUT
}

//IsInvalidNode returns true is node/cluster is down
func IsInvalidNode(err error) bool {
	if err == nil {
		return false
	}
	if err == NodeDownError {
		return true
	}
	aeroError, ok := err.(types.AerospikeError)
	if !ok {
		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
		if aeroError, ok = err.(types.AerospikeError); !ok {
			return strings.Contains(err.Error(), "connection refused")
		}

	}
	return aeroError.ResultCode() == types.INVALID_NODE_ERROR
}
