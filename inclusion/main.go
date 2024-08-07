package main

import (
	"fmt"

	"github.com/datatrails/go-datatrails-common/azblob"
	"github.com/datatrails/go-datatrails-logverification/logverification"
)

/**
 * Tests the inclusion of a datatrails event.
 *
 * This is achieved by creating an inclusion proof of the event,
 *  then verifying the inclusion proof.
 *
 * This proves that the datatrails event is included in the merklelog.
 */

const (

	// event is the event in json format as returned by the datatrails events API
	//
	// the event used in this demo is from a breast cancer diagnosing AI model sample
	//  found here: https://app.datatrails.ai/archivist/publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/71d7ab65-359b-40d9-9bbd-102ec2092601
	event = `
	{
		"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/71d7ab65-359b-40d9-9bbd-102ec2092601",
		"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
		"event_attributes": {
			"arc_description": "Approving Model",
			"arc_display_type": "Model Approval",
			"approvers": "Product Team"
		},
		"asset_attributes": {
			"model_version": "mcbdc.01.0.0",
			"modelcard_version": "2.0.0",
			"datacard_version": "2.0.0"
		},
		"operation": "Record",
		"behaviour": "RecordEvidence",
		"timestamp_declared": "2024-05-07T20:32:00Z",
		"timestamp_accepted": "2024-05-07T20:32:00Z",
		"timestamp_committed": "2024-05-07T20:32:27.235Z",
		"principal_declared": {
			"issuer": "",
			"subject": "",
			"display_name": "",
			"email": ""
		},
		"principal_accepted": {
			"issuer": "",
			"subject": "",
			"display_name": "",
			"email": ""
		},
		"confirmation_status": "CONFIRMED",
		"transaction_id": "0x224d41c6d984cb67e52274d62a48cd31fce39b6731f25f827d4c59a9cfdff427",
		"block_number": 7030,
		"transaction_index": 0,
		"from": "0x344b47d0FC35a551bd8a7Db4999226C04E764db3",
		"tenant_identity": "tenant/f023005c-000f-4a57-b2fe-eef425f243ad",
		"merklelog_entry": {
			"commit": {
				"index": "499",
				"idtimestamp": "018f54c1f0640dca00"
			},
			"confirm": {
				"mmr_size": "501",
				"root": "AsPmdY7mI1E4Hpkut1e1dYhj+gsRBS2c4NNLvZ0NMBg=",
				"timestamp": "1715113947353",
				"idtimestamp": "",
				"signed_tree_head": ""
			},
			"unequivocal": null
		}
	}
	`

	publicTenantID = "tenant/6ea5cd00-c711-3649-6914-7b125928bbb4"

	// merklelog reader configuration
	container = "merklelogs"
	url       = "https://app.datatrails.ai/verifiabledata"
)

// InclusionDemo of a public datatrails event
func InclusionDemo(eventJson []byte) (verified bool, err error) {

	// then create the merklelog reader
	reader, err := azblob.NewReaderNoAuth(url, azblob.WithContainer(container))
	if err != nil {
		return false, err
	}

	verifiableEvent, err := logverification.NewVerifiableEvent(eventJson)
	if err != nil {
		return false, err
	}

	// now verify the public event is in the merklelog
	return logverification.VerifyEvent(reader, *verifiableEvent, logverification.WithMassifTenantId(publicTenantID))

}

// Demo of the inclusion of a public datatrails event
func main() {

	verified, err := InclusionDemo([]byte(event))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
	}

	fmt.Printf("\nEvent included on merkle log: %v\n", verified)

}
