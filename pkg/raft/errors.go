package raft

import "errors"

var (
	ErrLogIndexNotContinuous = errors.New("log index is not continuous")
	ErrStopped               = errors.New("raft is stopped")
)
