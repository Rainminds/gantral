package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// HashContext computes a deterministic SHA-256 hash of a JSON-compatible map.
// It ensures keys are sorted to guarantee the same hash for the same content.
func HashContext(data map[string]interface{}) (string, error) {
	// Handle nil or empty map consistently as empty JSON object "{}"
	if len(data) == 0 {
		hash := sha256.Sum256([]byte("{}"))
		return hex.EncodeToString(hash[:]), nil
	}

	// 1. Canonicalize by marshalling keys in sorted order.
	// We rely on standard encoding/json which sorts map keys (RFC 8785 behavior).
	// WARNING: Arrays/Slices are ORDERED sequences. We DO NOT sort them.
	// If the context contains Sets (unordered lists), the CALLER must sort them
	// to ensure deterministic hashing. Auto-sorting here would corrupt valid sequences (e.g. history logs).

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// 2. Hash
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}
