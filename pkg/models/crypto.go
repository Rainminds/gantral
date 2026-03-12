package models

// Signer defines the interface for cryptographically signing authoritative artifacts.
// TRD A.4 explicitly mandates that the control plane signs the artifact hash.
type Signer interface {
	// Sign generates a cryptographic signature over the provided hash payload.
	// It returns the signature bytes and the algorithm identifier used (e.g. "EdDSA-Ed25519").
	Sign(hash []byte) (signature []byte, algorithm string, err error)
}

// SignatureVerifier defines the interface for verifying artifact signatures
// during replay.
type SignatureVerifier interface {
	// Verify checks if the given signature is valid for the provided hash using the specified algorithm.
	Verify(hash, signature []byte, algorithm string) error
}

// TimeStamper defines the interface for generating cryptographic timestamp tokens
// to prove existence time
type TimeStamper interface {
	// Token generates a timestamp token binding the given hash to a specific time.
	// It returns the token bytes and the algorithm identifier (e.g. "RFC3161").
	Token(hash []byte) (token []byte, algorithm string, err error)
}

// TimestampVerifier defines the interface for verifying timestamp tokens
// during replay.
type TimestampVerifier interface {
	// Verify checks if the timestamp token is cryptographically valid for the given hash.
	Verify(hash, token []byte, algorithm string) error
}
