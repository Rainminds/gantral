package models

import (
	"testing"
)

func TestCanonicalSerialization_Ordering(t *testing.T) {
	// Two maps with same keys inserted in different orders
	m1 := map[string]interface{}{
		"z": 1,
		"a": "text",
		"m": true,
	}
	m2 := map[string]interface{}{
		"a": "text",
		"m": true,
		"z": 1,
	}

	b1, err := CanonicalSerialize(m1)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := CanonicalSerialize(m2)
	if err != nil {
		t.Fatal(err)
	}

	if string(b1) != string(b2) {
		t.Errorf("Serialization differs: %s vs %s", b1, b2)
	}
	expected := `{"a":"text","m":true,"z":1}`
	if string(b1) != expected {
		t.Errorf("Expected %s, got %s", expected, b1)
	}
}

func TestCanonicalSerialization_NullsAndHtml(t *testing.T) {
	m := map[string]interface{}{
		"optional": nil,
		"html":     "<script>&",
	}
	b, err := CanonicalSerialize(m)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"html":"<script>&","optional":null}`
	if string(b) != expected {
		t.Errorf("Expected %s, got %s", expected, b)
	}
}

func TestCanonicalSerialization_Floats(t *testing.T) {
	m := map[string]interface{}{
		"float": 1.2345,
	}
	b, err := CanonicalSerialize(m)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"float":1.2345}`
	if string(b) != expected {
		t.Errorf("Expected %s, got %s", expected, b)
	}
}
