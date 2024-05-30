package logverification

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"sort"

	"github.com/datatrails/go-datatrails-common-api-gen/assets/v2/assets"
	"github.com/datatrails/go-datatrails-merklelog/massifs"
	"github.com/datatrails/go-datatrails-simplehash/simplehash"
	"google.golang.org/protobuf/encoding/protojson"
)

// EventDetails contains key information for verifying inclusion of merkle log events
type EventDetails struct {
	eventID     string
	tenantID    string
	eventHash   []byte
	merkleLog   *assets.MerkleLogEntry
	massifIndex uint64
}

// HashEvent creates a hash of the supplied event in the canonical V3 format.
func HashEvent(idBytes []byte, err error, v3event simplehash.V3Event) ([]byte, error) {
	hasher := sha256.New()

	domainSeparator := []byte{byte(LeafTypePlain)}
	hasher.Write(domainSeparator)

	hasher.Write(idBytes)

	err = simplehash.V3HashEvent(hasher, v3event)
	if err != nil {
		return nil, err
	}

	hash := hasher.Sum(nil)

	return hash, nil
}

// getIdTimestamp retrieves the ID Timestamp from a wrapped event
func getIdTimestamp(timestamp string) ([]byte, error) {

	// Note that we store the idtimestamp in the event data with the epoch
	// prefixed to it. So that no matter when the event data is examined, that
	// information is available. We do NOT include the epoch in the hash because
	// we will never alow a single log to span multiple epochs. We will just
	// start a new log. And we preserve the unified history provability by
	// puting a commitment to the final state of the previous log as the first
	// entry in the new log.

	id, epoch, err := massifs.SplitIDTimestampHex(timestamp)
	if err != nil {
		return nil, err
	}
	if epoch < 1 {
		return nil, errors.New("epoch is before datatrails existed")
	}
	if epoch > 1 {
		return nil, errors.New("epoch is after Jan 2038")
	}
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return b, nil
}

// ParseEventList takes a json list of events returned by the datatrails events API
//
//	and returns an mmrIndex ascending, sorted list of golang list of event details whose members are easier to access.
func ParseEventList(eventsJson []byte) ([]EventDetails, error) {

	// get the event list out of events
	eventListJson := struct {
		Events []json.RawMessage `json:"events"`
	}{}
	err := json.Unmarshal(eventsJson, &eventListJson)
	if err != nil {
		return nil, err
	}

	events := []EventDetails{}
	for _, eventJson := range eventListJson.Events {

		// special care is needed here to deal with uint64 types. json marshal /
		// un marshal treats them as strings because they don't fit in a
		// javascript Number

		// Unmarshal into a generic type to get just the bits we need. Use
		// defered decoding to get the raw merklelog entry as it must be
		// unmarshaled using protojson and the specific generated target type.
		entry := struct {
			Identity       string `json:"identity,omitempty"`
			TenantIdentity string `json:"tenant_identity,omitempty"`
			// Note: the proof_details top level field can be ignored here because it is a 'oneof'
			MerklelogEntry json.RawMessage `json:"merklelog_entry,omitempty"`
		}{}
		err := json.Unmarshal(eventJson, &entry)
		if err != nil {
			return nil, err
		}

		merkleLog := &assets.MerkleLogEntry{}
		err = protojson.Unmarshal(entry.MerklelogEntry, merkleLog)
		if err != nil {
			return nil, err
		}

		// Hash the events we want to verify using the v3 schema
		v3event, err := simplehash.V3FromEventJSON(eventJson)
		if err != nil {
			return nil, err
		}

		idBytes, err := getIdTimestamp(merkleLog.Commit.Idtimestamp)
		if err != nil {
			return nil, err
		}

		eventHash, err := HashEvent(idBytes, err, v3event)
		if err != nil {
			return nil, err
		}

		eventDetails := EventDetails{
			eventID:   entry.Identity,
			tenantID:  entry.TenantIdentity,
			eventHash: eventHash,
			merkleLog: merkleLog,
		}
		events = append(events, eventDetails)
	}

	// Sorting the events by MMR index implies that they're sorted in log append order.
	sort.Slice(events, func(i, j int) bool {
		return events[i].merkleLog.Commit.Index < events[j].merkleLog.Commit.Index
	})

	return events, nil

}
