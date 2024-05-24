package main

import (
	"testing"

	"github.com/datatrails/go-datatrails-common/logger"
	"github.com/stretchr/testify/assert"
)

// TestIntegrityDemo tests the sample public event
//
//	is included on the merklelog.
func TestIntegrityDemo(t *testing.T) {

	// TODO: remove logging in azblob package, so we don't need a logger
	logger.New("NOOP")
	defer logger.OnExit()

	verified, err := IntegrityDemo([]byte(event))

	assert.Equal(t, nil, err)
	assert.Equal(t, true, verified)

}
