package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIntegrityDemo tests the sample public event
//
//	is included on the merklelog.
func TestIntegrityDemo(t *testing.T) {

	verified, err := IntegrityDemo([]byte(event))

	assert.Equal(t, nil, err)
	assert.Equal(t, true, verified)

}
