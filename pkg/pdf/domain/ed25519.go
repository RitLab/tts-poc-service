package domain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"tts-poc-service/models"
)

type Ed25519Key struct {
	Key string
}

type Ed25519BlockKey struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

type Ed25519Signature struct {
	Signature []byte
}

func LoadKey(b *models.Block) *Ed25519BlockKey {
	privKey, _ := hex.DecodeString(b.PrivateKey)
	pubKey, _ := hex.DecodeString(b.PublicKey)
	privateKey := ed25519.PrivateKey(privKey)
	publicKey := ed25519.PublicKey(pubKey)
	return &Ed25519BlockKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func (e *Ed25519Key) GeneratePublicAndPrivateKey() (*Ed25519BlockKey, error) {
	// Derive a 32-byte private key from the string using SHA-256
	hash := sha256.Sum256([]byte(e.Key))
	privateKey := ed25519.NewKeyFromSeed(hash[:])

	// Extract the public key from the private key
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return &Ed25519BlockKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

func (e *Ed25519BlockKey) SignData(data []byte) []byte {
	return ed25519.Sign(e.PrivateKey, data)
}

func (e *Ed25519BlockKey) VerifySignature(signature, data []byte) bool {
	return ed25519.Verify(e.PublicKey, data, signature)
}
