package server

import "errors"

var (
	// ErrHandlerConfigIsNotFull will be returned if some handlers Config fields are not populated
	ErrHandlerConfigIsNotFull = errors.New("some handler Config fields are not populated")
	// ErrServerConfigIsNotFull will be returned if some servers Config fields are not populated
	ErrServerConfigIsNotFull = errors.New("some server Config fields are not populated")
	// ErrRouterConfigIsNotFull will be returned if some servers Config fields are not populated
	ErrRouterConfigIsNotFull = errors.New("some router Config fields are not populated")
)
