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

const (

	// eventList is the event list in json format as returned by the datatrails events API
	//
	// the events used in this demo are from a breast cancer diagnosing AI model sample
	//  found here: https://app.datatrails.ai/archivist/publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6
	eventList = `
	{
		"events": [
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/9a192afe-9253-44d7-8585-c48f237f2134",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"date": "Tue May  7 15:33:28 2024",
					"reason": "Model is mispredicting, causing incorrect diagnosis",
					"version": "mcbdc.01.0.0",
					"arc_description": "Model Deactivated",
					"arc_display_type": "Model Deactivation"
				},
				"asset_attributes": {
					"deactivated": "Y"
				},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:33:28Z",
				"timestamp_accepted": "2024-05-07T20:33:28Z",
				"timestamp_committed": "2024-05-07T20:33:55.826Z",
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
						"index": "511",
						"idtimestamp": "018f54c34a730ce300"
					},
					"confirm": {
						"mmr_size": "512",
						"root": "SnBhDOt7lF/aTK48db1qk0/86dluDr+ypLqGWLy2/5g=",
						"timestamp": "1715114035953",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/e7a369a9-06be-4edd-9bb5-17bc4c75d095",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"comments": "model is mispredicting, something is wrong",
					"stable": "N",
					"arc_description": "Unstable Model Operation",
					"arc_display_type": "Model Operation"
				},
				"asset_attributes": {},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:33:25Z",
				"timestamp_accepted": "2024-05-07T20:33:25Z",
				"timestamp_committed": "2024-05-07T20:33:26.461Z",
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
						"index": "502",
						"idtimestamp": "018f54c2d7be0dca00"
					},
					"confirm": {
						"mmr_size": "511",
						"root": "f6/H7cQ0Ilr//BmwWC76KnGwai0DU1g1bfClLSJWwjU=",
						"timestamp": "1715114006735",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/df0c9067-8058-4b0f-9c87-61a90ed8ad11",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"deployed": "Wisconsin State Hospital East, Wisconsin State Hospital West",
					"release_date": "Tue May  7 15:32:45 2024",
					"arc_description": "Releasing Model",
					"arc_display_type": "Model Release"
				},
				"asset_attributes": {},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:32:46Z",
				"timestamp_accepted": "2024-05-07T20:32:46Z",
				"timestamp_committed": "2024-05-07T20:33:13.089Z",
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
						"index": "501",
						"idtimestamp": "018f54c2a3820dca00"
					},
					"confirm": {
						"mmr_size": "502",
						"root": "pVNHzmEjjE5c+PPdApxVtu7QmqgkoQYIWa72V+FsYX0=",
						"timestamp": "1715113993205",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/71d7ab65-359b-40d9-9bbd-102ec2092601",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"approvers": "Product Team",
					"arc_description": "Approving Model",
					"arc_display_type": "Model Approval"
				},
				"asset_attributes": {
					"datacard_version": "2.0.0",
					"model_version": "mcbdc.01.0.0",
					"modelcard_version": "2.0.0"
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
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/7f6fb741-c73b-49a7-8784-5c359d73d993",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"arc_description": "Testing Model",
					"arc_display_type": "Model Testing",
					"comments": "Model was trained on dataset v3.0.0",
					"date": "Tue May  7 15:31:58 2024",
					"features": "Tested all features listed in datacard",
					"model_version": "mcbdc.01.0.0",
					"results": "Pass",
					"testers": "RAI QA Team"
				},
				"asset_attributes": {},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:31:58Z",
				"timestamp_accepted": "2024-05-07T20:31:58Z",
				"timestamp_committed": "2024-05-07T20:31:59.600Z",
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
						"index": "498",
						"idtimestamp": "018f54c184710ce300"
					},
					"confirm": {
						"mmr_size": "499",
						"root": "ZLfenOkN9Wy2+bljplPMvfHhu4vUhOelXf7slmXQm3U=",
						"timestamp": "1715113919750",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/d55d6843-6f8a-449d-862b-99126f186c1b",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"arc_description": "Publishing Model Card",
					"arc_display_type": "Publish",
					"document_version_authors": [
						{
							"display_name": "Development Team",
							"email": "raiteam@demo.com"
						}
					]
				},
				"asset_attributes": {
					"document_status": "Published",
					"document_version": "2.0.0",
					"document_document": {
						"arc_file_name": "BCD_ModelCard_v2.pdf",
						"arc_attribute_type": "arc_attachment",
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "f84b351f15f942fd8ab53a959da446112f4c04b919942ea378ca1e6981dbf71e",
						"arc_blob_identity": "blobs/d53b9530-85af-4ce2-b88a-14e3c1381f8e",
						"arc_display_name": "published_document"
					},
					"document_hash_alg": "SHA256",
					"document_hash_value": "f84b351f15f942fd8ab53a959da446112f4c04b919942ea378ca1e6981dbf71e"
				},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:31:56Z",
				"timestamp_accepted": "2024-05-07T20:31:56Z",
				"timestamp_committed": "2024-05-07T20:31:57.758Z",
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
						"index": "495",
						"idtimestamp": "018f54c17d3f0ce300"
					},
					"confirm": {
						"mmr_size": "498",
						"root": "sTnRzIxhhHptTzMOqGJHeLuf3Ukrz1L7tLzeUCF4vWg=",
						"timestamp": "1715113917874",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/443f850f-294d-41c4-9c62-7e44d5344330",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"arc_display_type": "Model Card Approval",
					"approvers": "Product Team",
					"arc_description": "Approving Model Card"
				},
				"asset_attributes": {
					"model_version": "0.0.0",
					"modelcard_version": "2.0.0",
					"datacard_version": "2.0.0"
				},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:31:55Z",
				"timestamp_accepted": "2024-05-07T20:31:55Z",
				"timestamp_committed": "2024-05-07T20:31:55.852Z",
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
						"index": "494",
						"idtimestamp": "018f54c175cd0dca00"
					},
					"confirm": {
						"mmr_size": "495",
						"root": "oTiPumw6HRFfp0GlmNOQj/cX2stK4WNdik1wiWIJYU8=",
						"timestamp": "1715113916026",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/87d9a1a5-5120-4613-bf40-6918297e093c",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"mc_version": "2.0.0",
					"model_card": {
						"arc_file_name": "BCD_ModelCard_v2.pdf",
						"arc_attribute_type": "arc_attachment",
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "f84b351f15f942fd8ab53a959da446112f4c04b919942ea378ca1e6981dbf71e",
						"arc_blob_identity": "blobs/ede0c293-0341-4def-a81a-739246680625",
						"arc_display_name": "model_card"
					},
					"reviewers": "RAI Stakeholders",
					"approved": "Y",
					"arc_description": "Reviewing Model Card - Pass",
					"arc_display_type": "Model Card Review",
					"comments": "Model Card looks good",
					"date": "Tue May  7 15:31:52 2024"
				},
				"asset_attributes": {},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:31:53Z",
				"timestamp_accepted": "2024-05-07T20:31:53Z",
				"timestamp_committed": "2024-05-07T20:31:54.150Z",
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
						"index": "492",
						"idtimestamp": "018f54c16f270dca00"
					},
					"confirm": {
						"mmr_size": "494",
						"root": "/VstI/yIg6uVFMwUyoQss21u3iUSxmCjg/tgbHtVOBA=",
						"timestamp": "1715113914269",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/2abc8706-b50a-47e2-88a3-64a70bbe5bce",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"features": "Tested all features listed in datacard",
					"model_version": "mcbdc.01.0.0",
					"results": "Initial metrics listed within model card",
					"testers": "RAI QA Team",
					"arc_description": "Testing Model",
					"arc_display_type": "Model Testing",
					"comments": "Model was trained on dataset v3.0.0",
					"date": "Tue May  7 15:31:08 2024"
				},
				"asset_attributes": {},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:31:08Z",
				"timestamp_accepted": "2024-05-07T20:31:08Z",
				"timestamp_committed": "2024-05-07T20:31:35.942Z",
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
						"index": "491",
						"idtimestamp": "018f54c128070ce300"
					},
					"confirm": {
						"mmr_size": "492",
						"root": "N7mYlGmoYukxgf6JzGO0UhF1JIUOJLCf8TxDC9C6wBk=",
						"timestamp": "1715113896062",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/40bdad1c-3a2c-4774-be67-e57102a19c99",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"arc_description": "Publishing Data Card",
					"arc_display_type": "Publish",
					"document_version_authors": [
						{
							"email": "raiteam@demo.com",
							"display_name": "Development Team"
						}
					]
				},
				"asset_attributes": {
					"document_hash_value": "5147787c57a0943e823e78416135aa89c1fa317a3f020e27a2a7b7719dc9a436",
					"document_status": "Published",
					"document_version": "2.0.0",
					"document_document": {
						"arc_attribute_type": "arc_attachment",
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "5147787c57a0943e823e78416135aa89c1fa317a3f020e27a2a7b7719dc9a436",
						"arc_blob_identity": "blobs/2460e045-1f1a-4aa5-a467-136f2491168d",
						"arc_display_name": "published_document",
						"arc_file_name": "BCD_DataCard_v2.pdf"
					},
					"document_hash_alg": "SHA256"
				},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:31:06Z",
				"timestamp_accepted": "2024-05-07T20:31:06Z",
				"timestamp_committed": "2024-05-07T20:31:07.313Z",
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
						"index": "487",
						"idtimestamp": "018f54c0b8320ce300"
					},
					"confirm": {
						"mmr_size": "491",
						"root": "w2TR6cbx1cRShOSQuEJkgOyYIipK+HCf/Nev8fkTt00=",
						"timestamp": "1715113867427",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/8fe3af12-4052-45e8-bb8c-7eba4cd04ead",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"arc_display_type": "Data Card Approval",
					"approvers": "Product Team",
					"arc_description": "Approving Data Card"
				},
				"asset_attributes": {
					"model_version": "0.0.0",
					"modelcard_version": "0.0.0",
					"datacard_version": "2.0.0"
				},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:30:08Z",
				"timestamp_accepted": "2024-05-07T20:30:08Z",
				"timestamp_committed": "2024-05-07T20:30:38.104Z",
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
						"index": "486",
						"idtimestamp": "018f54c046190ce300"
					},
					"confirm": {
						"mmr_size": "487",
						"root": "BuArZXNf+rAAqr1FBrF9GyCJrRfVerW6JisdAAOdbe0=",
						"timestamp": "1715113838217",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/3b60a93e-65a7-4d12-96ea-70ef8d85b24f",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"dataset": {
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "5594fb3ae1f345a03796c3678c84164d812379f326abf3e1c8942df2f8f610bf",
						"arc_blob_identity": "blobs/b14f6c25-cf9a-475b-ba47-d4e8f9957a8b",
						"arc_display_name": "data_set",
						"arc_file_name": "breast-cancer_v3.data",
						"arc_attribute_type": "arc_attachment"
					},
					"arc_description": "Reviewing Data Card - Pass",
					"data_version": "2.0.0",
					"approved": "Y",
					"arc_display_type": "Data Card Review",
					"data_card": {
						"arc_blob_identity": "blobs/368966c7-e196-4f66-8f7b-8fdfb5d45b91",
						"arc_display_name": "data_card",
						"arc_file_name": "BCD_DataCard_v2.pdf",
						"arc_attribute_type": "arc_attachment",
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "5147787c57a0943e823e78416135aa89c1fa317a3f020e27a2a7b7719dc9a436"
					},
					"date": "Tue May  7 15:29:09 2024",
					"reviewers": "RAI Stakeholders",
					"comments": "Data Card and Dataset looks good"
				},
				"asset_attributes": {},
				"operation": "Record",
				"behaviour": "RecordEvidence",
				"timestamp_declared": "2024-05-07T20:29:10Z",
				"timestamp_accepted": "2024-05-07T20:29:10Z",
				"timestamp_committed": "2024-05-07T20:29:43.545Z",
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
						"index": "484",
						"idtimestamp": "018f54bf70fa0dca00"
					},
					"confirm": {
						"mmr_size": "486",
						"root": "GDHZZAdlYwJIger/lt2KcW+/WsQYB59W/y1RXHTYLFo=",
						"timestamp": "1715113783660",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			},
			{
				"identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6/events/99c77eae-d1d7-4210-b55f-e2be638b974b",
				"asset_identity": "publicassets/3ea5aca3-da02-4bae-b6d0-85a5ab586ed6",
				"event_attributes": {
					"arc_access_policy_always_read": [
						{
							"wallet": "0x344b47d0FC35a551bd8a7Db4999226C04E764db3",
							"tessera": "SmL4PHAHXLdpkj/c6Xs+2br+hxqLmhcRk75Hkj5DyEQ="
						}
					],
					"arc_access_policy_asset_attributes_read": [
						{
							"0x4609ea6bbe85F61bc64760273ce6D89A632B569f": "wallet",
							"SmL4PHAHXLdpkj/c6Xs+2br+hxqLmhcRk75Hkj5DyEQ=": "tessera",
							"attribute": "*"
						}
					],
					"arc_access_policy_event_arc_display_type_read": [
						{
							"0x4609ea6bbe85F61bc64760273ce6D89A632B569f": "wallet",
							"SmL4PHAHXLdpkj/c6Xs+2br+hxqLmhcRk75Hkj5DyEQ=": "tessera",
							"value": "*"
						}
					]
				},
				"asset_attributes": {
					"document_version": "0.0.0",
					"id": "mw0wv4",
					"document_hash_alg": "SHA256",
					"vendor": "AIVendor_B",
					"arc_description": "Breast Cancer Detection - Wisconsin",
					"arc_display_type": "AI Item",
					"arc_display_name": "Breast Cancer Detection",
					"document_hash_value": "52e8cefae041d3a7ca9efd08538b61976c64b41c43f788bef1e8979c96b1cdff",
					"arc_primary_image": {
						"arc_attribute_type": "arc_attachment",
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "52e8cefae041d3a7ca9efd08538b61976c64b41c43f788bef1e8979c96b1cdff",
						"arc_blob_identity": "blobs/f9e458b5-3b34-43fe-9546-49f952a607fc",
						"arc_display_name": "arc_primary_image",
						"arc_file_name": "pexels-jason-deines-6246853.jpg"
					},
					"document_document": {
						"arc_attribute_type": "arc_attachment",
						"arc_blob_hash_alg": "SHA256",
						"arc_blob_hash_value": "52e8cefae041d3a7ca9efd08538b61976c64b41c43f788bef1e8979c96b1cdff",
						"arc_blob_identity": "blobs/c2b4fdc9-d8a6-4bea-9402-0c1ff6f93695",
						"arc_display_name": "ai_image",
						"arc_file_name": "pexels-jason-deines-6246853.jpg"
					},
					"arc_profile": "Document"
				},
				"operation": "NewAsset",
				"behaviour": "AssetCreator",
				"timestamp_declared": "2024-05-07T20:29:05Z",
				"timestamp_accepted": "2024-05-07T20:29:05Z",
				"timestamp_committed": "2024-05-07T20:29:07Z",
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
						"index": "483",
						"idtimestamp": "018f54bee2390ce300"
					},
					"confirm": {
						"mmr_size": "484",
						"root": "/zdItj0xBrjeWsu2FoL5v1+YurBLmfWJTYylu46wRNs=",
						"timestamp": "1715113747319",
						"idtimestamp": "",
						"signed_tree_head": ""
					},
					"unequivocal": null
				}
			}
		],
		"next_page_token": ""
	}
	`

	publicTenantID = "tenant/6ea5cd00-c711-3649-6914-7b125928bbb4"

	// merklelog reader configuration
	container = "merklelogs"
	url       = "https://app.datatrails.ai/verifiabledata"
)

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
		fmt.Printf("\nFailed verification: %v\n", err)
		os.Exit(1)
	}

	// if we have any omitted events then the verification fails.
	//  an omitted event is an event on the merklelog that is NOT
	//  included in the given list of events.
	if len(omittedEvents) > 0 {
		fmt.Printf("\nFailed verification, omitted events: %v", omittedEvents)
		os.Exit(1)
	}

	fmt.Println("Complete List of events included on merkle log")

}
