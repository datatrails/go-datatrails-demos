package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/datatrails/go-datatrails-common/azblob"
	"github.com/datatrails/go-datatrails-demos/logverification"
	"github.com/datatrails/go-datatrails-merklelog/massifs"
)

func ConsistencyDemo() (verified bool, err error) {

	// First we need to get the existing signed log state from a trusted source, one we saved earlier.
	//
	// The log state saved is from the massif that contains an event for the breast cancer diagnosing AI model sample.
	// The event can be found here: https://app.datatrails.ai/archivist/publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/9a192afe-9253-44d7-8585-c48f237f2134
	//
	// This does the following:
	//  1. gets the saved signed log state
	//  2. verifies the signature of the signed log state
	//  3. unmarshals the signed log state into a golang data structure.
	existingLogState, err := ExistingSignedState()
	if err != nil {
		return false, err
	}

	// Now we get a future log state and confirm that our earlier event continues to be consistently recorded.
	//  the existing signed log state.
	//
	// The event we are basing the newer log state on is:
	//   https://app.datatrails.ai/archivist/publicassets/fe022486-3272-4d44-aab5-765a37c17b85/events/3e7a16dd-01d6-44f5-870d-abb9c56d154b

	// Get the signed state for the newer log state
	hasher := sha256.New()
	codec, err := massifs.NewRootSignerCodec()
	if err != nil {
		return false, err
	}

	reader, err := azblob.NewReaderNoAuth(url, azblob.WithContainer(container))
	if err != nil {
		return false, err
	}

	massifIndex, err := massifs.MassifIndexFromMMRIndex(logverification.DefaultMassifHeight, newStateMMRIndex)
	if err != nil {
		return false, err
	}
	signedState, err := logverification.SignedLogState(
		context.Background(), reader, hasher, codec, publicTenantID,
		massifIndex,
	)
	if err != nil {
		return false, err
	}

	// Now verify the signed state using the datatrails seal verification key
	verificationKey, err := VerificationKeyFromFile()
	if err != nil {
		return false, err
	}

	err = signedState.VerifyWithPublicKey(verificationKey, nil)
	if err != nil {
		return false, err
	}

	// unmarshal the signed log state into a golang data structure.
	logState, err := logverification.LogState(signedState, codec)
	if err != nil {
		return false, err
	}

	// Now we have 2 log states that we have verified the signature using the datatrails seal verification key.
	//
	// The first log state is taken from when an event for the breast cancer diagnosing AI model sample is on the log.
	// The event can be found here: https://app.datatrails.ai/archivist/publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/9a192afe-9253-44d7-8585-c48f237f2134
	//
	// The second log state is taken from a later point in time when the above event, as well as the following event is on the log:
	// https://app.datatrails.ai/archivist/publicassets/fe022486-3272-4d44-aab5-765a37c17b85/events/3e7a16dd-01d6-44f5-870d-abb9c56d154b
	//
	// We want to make sure that the second log state is appended from the first

	return logverification.VerifyConsistency(context.Background(), hasher, reader, publicTenantID, existingLogState, logState)

}

func main() {

	verified, err := ConsistencyDemo()

	if err != nil {
		fmt.Printf("Failed to verify the consistency of the two log states: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Two log state verification consistency is: %v\n", verified)
}
