package main

import (
	"fmt"
	"os"

	"github.com/datatrails/go-datatrails-common/azblob"
	"github.com/datatrails/go-datatrails-logverification/logverification"
)

/**
 * Tests the integrity of a datatrails event.
 *
 * This is achieved by creating an inclusion proof of the event,
 *  then verifying the inclusion proof.
 *
 * This proves that the datatrails event is included in the merklelog.
 */

// CompletenessDemo of a list of public datatrails events
func CompletenessDemo(eventsJson []byte) (omittedEvents []uint64, err error) {

	// then create the merklelog reader
	reader, err := azblob.NewReaderNoAuth(url, azblob.WithContainer(container))
	if err != nil {
		return nil, err
	}

	verifiableEvents, err := logverification.NewVerifiableEvents(eventsJson)
	if err != nil {
		return nil, err
	}

	// now verify the public event is in the merklelog
	return logverification.VerifyList(reader, verifiableEvents, logverification.WithTenantId(publicTenantID))

}

// Demo of the completeness of a public datatrails event
func main() {

	omittedEvents, err := CompletenessDemo([]byte(eventList))
	if err != nil {
		fmt.Printf("\nFailed Complete List verification: %v\n", err)
		os.Exit(1)
	}

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
	if len(omittedEvents) > 0 {
		fmt.Printf("\nFailed Complete List verification, omitted events mmrIndexs: %v\n", omittedEvents)
		os.Exit(1)
	}

	fmt.Println("Complete List of events included on merkle log")

}
