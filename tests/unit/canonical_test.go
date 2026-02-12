package unit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/Rainminds/gantral/internal/artifact"
)

// Test_Canonicalization_KeyOrder_SectionK verifies that map key order
// does not affect the generated hash.
// "Map key order invariance (top-level and nested)."
func Test_Canonicalization_KeyOrder_SectionK(t *testing.T) {
	t.Parallel()

	// Two maps with IDENTICAL content but inserted in different order?
	// Go maps are unordered by definition, but literal initialization might suggest order to humans.
	// Go's `json.Marshal` enforces sorting of keys.

	// Case 1: Simple Map
	m1 := map[string]interface{}{
		"a": 1,
		"b": 2,
	}
	m2 := map[string]interface{}{
		"b": 2,
		"a": 1,
	}

	h1, err := artifact.HashContext(m1)
	if err != nil {
		t.Fatalf("HashContext m1 failed: %v", err)
	}
	h2, err := artifact.HashContext(m2)
	if err != nil {
		t.Fatalf("HashContext m2 failed: %v", err)
	}

	if h1 != h2 {
		t.Errorf("Hash mismatch for identical maps!\nH1: %s\nH2: %s", h1, h2)
	}

	// Case 2: Nested Map
	nested1 := map[string]interface{}{
		"root": map[string]interface{}{
			"x": "foo",
			"y": "bar",
		},
	}
	nested2 := map[string]interface{}{
		"root": map[string]interface{}{
			"y": "bar",
			"x": "foo",
		},
	}

	hn1, _ := artifact.HashContext(nested1)
	hn2, _ := artifact.HashContext(nested2)

	if hn1 != hn2 {
		t.Errorf("Nested hash mismatch!\nH1: %s\nH2: %s", hn1, hn2)
	}
}

// Test_Canonicalization_NumericEdges_SectionK verifies numeric handling.
// "Numeric edge handling (1 vs 1.0)"
func Test_Canonicalization_NumericEdges_SectionK(t *testing.T) {
	t.Parallel()

	// In Go, `interface{}` validation depends on how json unmarshals it or how it's created.
	// 1 vs 1.0.
	// If created in Go as int vs float64, they serialize differently by default (1 vs 1).
	// json.Marshal(1) -> "1"
	// json.Marshal(1.0) -> "1" (Go trims suffix for whole floats)
	// json.Marshal(1.1) -> "1.1"

	mInt := map[string]interface{}{"val": 1}
	mFloat := map[string]interface{}{"val": 1.0}

	bInt, _ := json.Marshal(mInt)
	bFloat, _ := json.Marshal(mFloat)

	// Assert that Go's default marshaling treats them identical (good for stability)
	if string(bInt) != string(bFloat) {
		t.Logf("Notice: Go marshals 1 and 1.0 differently? Int: %s, Float: %s", bInt, bFloat)
		// If they differ, our hash will differ.
		// For Tier 1, we want STABLE bytes.
		// If the input implies 1 (integer), and another system sends 1.0 (float),
		// strict canonicalization usually demands a specific format (e.g. all numbers are floats, or preserve raw).
		// However, for pure Go map[string]interface{}, 1 and 1.0 usually marshal to "1".
	}

	hInt, _ := artifact.HashContext(mInt)
	hFloat, _ := artifact.HashContext(mFloat)

	if hInt != hFloat {
		// If this fails, it means we have ambiguity.
		// For strict enforcement, we might want to REJECT ambiguity or NORMALIZE.
		// The requirement says: "Numeric edge handling (1 vs 1.0; -0; NaN/Infinity rejected if representable in inputs)."
		// Go's standard lib treats 1 and 1.0 as "1".
		t.Errorf("Hash mismatch for 1 vs 1.0. This implies canonicalization instability across types.\nHInt: %s\nHFloat: %s", hInt, hFloat)
	}
}

// Test_Crypto_SectionM verifies correct SHA-256 usage.
func Test_Crypto_SectionM(t *testing.T) {
	t.Parallel()

	// Ensure HashContext uses SHA256
	input := map[string]interface{}{"test": "value"}
	hash, err := artifact.HashContext(input)
	if err != nil {
		t.Fatal(err)
	}

	// Verify length (SHA256 hex is 64 chars)
	if len(hash) != 64 {
		t.Errorf("Expected 64-char hex hash (SHA-256), got len %d: %s", len(hash), hash)
	}

	// Verify correctness against standard lib
	bytes, _ := json.Marshal(input)
	sum := sha256.Sum256(bytes)
	expected := hex.EncodeToString(sum[:])

	if hash != expected {
		t.Errorf("Hash mismatch against direct SHA256.\nFunc: %s\nExp:  %s", hash, expected)
	}
}
