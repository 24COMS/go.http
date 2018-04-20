package server

import "errors"

var (
	// ErrHandlerConfigIsNotFull will be returned if some handlers Config fields are not populated
	ErrHandlerConfigIsNotFull = errors.New("some handlers Config fields are not populated")
	// ErrServerConfigIsNotFull will be returned if some servers Config fields are not populated
	ErrServerConfigIsNotFull = errors.New("some servers Config fields are not populated")
)
