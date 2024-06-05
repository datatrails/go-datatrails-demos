package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	verified, err := ConsistencyDemo()

	assert.Equal(t, nil, err)
	assert.Equal(t, true, verified)

}
