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

	// If we have any omitted events then the verification fails.
	//
	// This is because we expected a complete list of events. Which
	//  means that ONLY events in the list are on the merklelog, within
	//  the range of the first event to the last event in the list.
	//
	// An omitted event is an event on the merklelog that is NOT
	//  included in the given list of events.
	//
	// NOTE: in other contexts, it is reasonable to have an in-complete
	//       list of events, where unrelated events are purposefully omitted.
	//
	assert.Equal(t, 0, len(omittedEvents))

}
