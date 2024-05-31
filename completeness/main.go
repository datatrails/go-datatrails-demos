package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/datatrails/go-datatrails-common/azblob"
	"github.com/datatrails/go-datatrails-demos/logverification"
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
func CompletenessDemo(eventJson []byte) (omittedEvents []uint64, err error) {

	// first we need to strip the 'public' prefix from all identity and assetIdentity fields
	//  as the hashing schema does not support public prefixing
	events := strings.ReplaceAll(eventList, "publicassets/", "assets/")

	// then create the merklelog reader
	reader, err := azblob.NewReaderNoAuth(url, azblob.WithContainer(container))
	if err != nil {
		return nil, err
	}

	// now verify the public event is in the merklelog
	return logverification.VerifyList(reader, []byte(events), logverification.WithTenantId(publicTenantID))

}

// Demo of the integrity of a public datatrails event
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
		fmt.Printf("\nFailed Complete List verification, omitted events: %v", omittedEvents)
		os.Exit(1)
	}

	fmt.Println("Complete List of events included on merkle log")

}
