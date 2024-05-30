package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/** TestCompletenessDemo tests the sample public events
 *   are included on the merklelog.
 *	 Also checks that the list is complete, i.e.
 *	 there are no events missing from the list, or have been added to the list
 *	 that are not on the merklelog.
 */
func TestCompletenessDemo(t *testing.T) {

	omittedEvents, err := CompletenessDemo([]byte(eventList))

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(omittedEvents))

}
