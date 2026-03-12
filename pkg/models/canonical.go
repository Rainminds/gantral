package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
)

// CanonicalSerialize serializes a payload map to a strictly deterministic,
// compact JSON byte array complying with TRD Appendix A.10.
func CanonicalSerialize(payload map[string]interface{}) ([]byte, error) {
	if payload == nil {
		return []byte("null"), nil
	}
	var buf bytes.Buffer
	if err := encodeMap(&buf, payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeValue(buf *bytes.Buffer, val interface{}) error {
	if val == nil {
		buf.WriteString("null")
		return nil
	}

	switch v := val.(type) {
	case string:
		// Encode string without HTML escaping to match generic JSON implementations
		var strBuf bytes.Buffer
		enc := json.NewEncoder(&strBuf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(v); err != nil {
			return err
		}
		res := strBuf.Bytes()
		// Remove the trailing newline added by Encoder
		if len(res) > 0 && res[len(res)-1] == '\n' {
			res = res[:len(res)-1]
		}
		buf.Write(res)
	case bool:
		if v {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case int:
		buf.WriteString(strconv.Itoa(v))
	case int32:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case float64:
		return encodeFloat(buf, v)
	case float32:
		return encodeFloat(buf, float64(v))
	case map[string]interface{}:
		return encodeMap(buf, v)
	case []interface{}:
		return encodeSlice(buf, v)
	case []Attestation:
		// Handle specific schema slices if any exist in the hashing payload
		slice := make([]interface{}, len(v))
		for i, x := range v {
			slice[i] = x
		}
		return encodeSlice(buf, slice)
	default:
		// Fallback for simple map structs or aliases
		return fmt.Errorf("canonical serialize: unsupported type %T", v)
	}
	return nil
}

func encodeFloat(buf *bytes.Buffer, f float64) error {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Errorf("NaN or Inf not allowed in canonical JSON")
	}
	// TRD A.10.7: Decimal base, no exponent notation
	buf.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
	return nil
}

func encodeMap(buf *bytes.Buffer, m map[string]interface{}) error {
	buf.WriteByte('{')
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// TRD A.10.2: Lexicographic ascending order
	sort.Strings(keys)

	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}

		var strBuf bytes.Buffer
		enc := json.NewEncoder(&strBuf)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(k)
		res := strBuf.Bytes()
		if len(res) > 0 && res[len(res)-1] == '\n' {
			res = res[:len(res)-1]
		}
		buf.Write(res)

		buf.WriteByte(':')

		if err := encodeValue(buf, m[k]); err != nil {
			return err
		}
	}
	buf.WriteByte('}')
	return nil
}

func encodeSlice(buf *bytes.Buffer, slice []interface{}) error {
	buf.WriteByte('[')
	for i, v := range slice {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := encodeValue(buf, v); err != nil {
			return err
		}
	}
	buf.WriteByte(']')
	return nil
}
