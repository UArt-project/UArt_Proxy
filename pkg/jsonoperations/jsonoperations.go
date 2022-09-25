// Package jsonoperations contains basic things needed to work with JSON.
package jsonoperations

import (
	"encoding/json"
	"fmt"
)

// Encode parses any data to JSON.
func Encode(data any) ([]byte, error) {
	encData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("v1.Encode: %w", err)
	}

	return encData, nil
}

// EncodeIndent parses any data to JSON with specified indentation.
func EncodeIndent(data any, prefix, indent string) ([]byte, error) {
	encData, err := json.MarshalIndent(data, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("v1.EncodeIndent: %w", err)
	}

	return encData, nil
}

// Decode parses any data to JSON.
func Decode(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("DecodeOrderData: %w", err)
	}

	return nil
}
