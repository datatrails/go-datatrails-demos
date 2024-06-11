package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsistencyDemo(t *testing.T) {

	verified, err := ConsistencyDemo()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, verified)

}
