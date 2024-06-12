package main

import (
	"github.com/datatrails/go-datatrails-common/cose"
	"github.com/datatrails/go-datatrails-logverification/logverification"
	"github.com/datatrails/go-datatrails-merklelog/massifs"
)

/**
 * Existing log state is an existing log state based on the breast cancer diagnosing AI model sample
 *  found here: https://app.datatrails.ai/archivist/publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/9a192afe-9253-44d7-8585-c48f237f2134
 *
 * The existing log state will take the root of the massif the event is found in, in this case the mmrIndex of 511 of the public tenant.
 */

// ExistingSignedState gets the existing signed state for the log at the massif where
//
//	the event for the breast cancer diagnosing AI model sample is found.
//	  The event can be found here: https://app.datatrails.ai/archivist/publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/9a192afe-9253-44d7-8585-c48f237f2134
//
// Then verifies the existing signed state signature against using the known veriication key.
func ExistingSignedState() (*massifs.MMRState, error) {
	signedState, err := cose.NewCoseSign1MessageFromCBOR(sampleSignedStateCbor)
	if err != nil {
		return nil, err
	}

	verificationKey, err := VerificationKeyFromFile()
	if err != nil {
		return nil, err
	}

	err = signedState.VerifyWithPublicKey(verificationKey, nil)
	if err != nil {
		return nil, err
	}

	codec, err := massifs.NewRootSignerCodec()
	if err != nil {
		return nil, err
	}

	return logverification.LogState(signedState, codec)
}
