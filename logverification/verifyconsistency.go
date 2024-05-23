package logverification

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"hash"

	"github.com/datatrails/go-datatrails-common/azblob"
	"github.com/datatrails/go-datatrails-merklelog/massifs"
	"github.com/datatrails/go-datatrails-merklelog/mmr"
)

// VerifyConsistency takes a massif context providing access to data from the past, and a massif
// context providing access to the current version of the log. It returns whether or not the
// new version of the log is consistent with the previous version (i.e. it contains all of the
// same nodes in the same positions.)
func VerifyConsistency(
	ctx context.Context,
	verificationKey ecdsa.PublicKey,
	hasher hash.Hash,
	blobReader azblob.Reader,
	massifContextBefore *massifs.MassifContext,
	massifContextNow *massifs.MassifContext,
	logStateNow *massifs.MMRState,
) (bool, error) {
	// Grab some core info about our backed up merkle log, which we'll need to prove consistency
	mmrSizeBefore := massifContextBefore.Count()
	rootBefore, err := mmr.GetRoot(mmrSizeBefore, massifContextBefore, hasher)
	if err != nil {
		return false, fmt.Errorf("VerifyConsistency failed: unable to get root for massifContextBefore: %w", err)
	}

	// We construct a proof of consistency between the backed up MMR log and the head of the log.
	consistencyProof, err := mmr.IndexConsistencyProof(mmrSizeBefore, logStateNow.MMRSize, massifContextNow, hasher)
	if err != nil {
		return false, errors.New("error")
	}

	// In order to verify the proof we take the hashes of all of the peaks in the backed up log.
	// The hash of each of these peaks guarantees the integrity of all of its child nodes, so we
	// don't need to check every hash.

	// Peaks returned as MMR positions (1-based), not MMR indices (0-based). The location of these
	// is deterministic: Given an MMR of a particular size, the peaks will always be in the same place.
	backupLogPeaks := mmr.Peaks(mmrSizeBefore)

	// Get the hashes of all of the peaks.
	backupLogPeakHashes, err := mmr.PeakBagRHS(massifContextNow, hasher, 0, backupLogPeaks)
	if err != nil {
		return false, errors.New("error")
	}

	// Lastly, verify the consistency proof using the peak hashes from our backed-up log. If this
	// returns true, then we can confidently say that everything in the backed-up log is in the state
	// of the log described by this signed state.
	verified := mmr.VerifyConsistency(hasher, backupLogPeakHashes, consistencyProof, rootBefore, logStateNow.Root)
	return verified, nil
}
