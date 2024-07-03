package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInclusionDemo tests the sample public event
//
//	is included on the merklelog.
func TestInclusionDemo(t *testing.T) {

	verified, err := InclusionDemo([]byte(event))

	assert.Equal(t, nil, err)
	assert.Equal(t, true, verified)

}
