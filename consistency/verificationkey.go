package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

/**
 * Verification key holds utilities for getting the public key from the pem file.
 */

const (
	publicVerificationKeyFile = "verificationkey.pem"
)

// VerificationKeyFromFile gets the datatrails public verification key used
//
//	to verify the signature of merklelog seals.
func VerificationKeyFromFile() (*ecdsa.PublicKey, error) {

	verificationKeyPem, err := os.ReadFile(publicVerificationKeyFile)
	if err != nil {
		return nil, err
	}

	verificationKeyPemblock, _ := pem.Decode(verificationKeyPem)
	parseResult, err := x509.ParsePKIXPublicKey(verificationKeyPemblock.Bytes)
	if err != nil {
		return nil, err
	}

	verificationKey := parseResult.(*ecdsa.PublicKey)

	return verificationKey, nil

}
