package logverification

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"hash"

	"github.com/datatrails/go-datatrails-common/azblob"
	"github.com/datatrails/go-datatrails-common/cbor"
	"github.com/datatrails/go-datatrails-common/logger"
	"github.com/datatrails/go-datatrails-merklelog/massifs"
	"github.com/datatrails/go-datatrails-merklelog/mmr"
)

// VerifySignature downloads the latest signed root from the given massif and verifies the signature.
func VerifySignature(
	ctx context.Context,
	tenantID string,
	verificationKey ecdsa.PublicKey,
	hasher hash.Hash,
	codec cbor.CBORCodec,
	reader azblob.Reader,
	massifIndex uint64,
) (*massifs.MMRState, error) {
	sealReader := massifs.NewSignedRootReader(logger.Sugar, reader, codec)

	// Fetch the latest signed state of the log
	signedStateNow, logStateNow, err := sealReader.GetLatestMassifSignedRoot(ctx, tenantID, uint32(massifIndex))
	if err != nil {
		return nil, fmt.Errorf("VerifyConsistency failed: unable to get latest signed root: %w", err)
	}

	massifReader := massifs.NewMassifReader(logger.Sugar, reader)
	massifContext, err := massifReader.GetMassif(ctx, tenantID, massifIndex)
	if err != nil {
		return nil, fmt.Errorf("VerifyConsistency failed: unable to get massif from storage for massif index: %v, err: %w",
			massifIndex, err)
	}

	// The log state at time of sealing is the Payload. It included the root, but this is removed
	// from the stored log state. This forces a verifier to recompute the merkle root from their view
	// of the data. If verification succeeds when this computed root is added to signedStateNow, then
	// we can be confident that DataTrails signed this state, and that the root matches your data.
	logStateNow.Root, err = mmr.GetRoot(logStateNow.MMRSize, &massifContext, hasher)
	if err != nil {
		return nil, fmt.Errorf("VerifyConsistency failed: unable to get root for massifContextNow: %w", err)
	}

	signedStateNow.Payload, err = codec.MarshalCBOR(logStateNow)
	if err != nil {
		return nil, fmt.Errorf("VerifyConsistency failed: unable to cbor encode log state: %w", err)
	}

	signatureVerificationError := signedStateNow.VerifyWithPublicKey(&verificationKey, nil)
	if signatureVerificationError != nil {
		return nil, fmt.Errorf("VerifyConsistency failed: signature verification failed: %w", err)
	}

	return &logStateNow, nil
}
