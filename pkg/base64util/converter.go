package base64util

import (
	"encoding/base64"
	"errors"
)

// Encoder provides base64 encoding and decoding functionality
type Encoder struct{}

// NewEncoder creates a new base64 encoder instance
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode converts text to base64
func (e *Encoder) Encode(text string) (string, error) {
	if text == "" {
		return "", errors.New("text cannot be empty")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	return encoded, nil
}

// Decode converts base64 text back to original text
func (e *Encoder) Decode(encodedText string) (string, error) {
	if encodedText == "" {
		return "", errors.New("encoded text cannot be empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", errors.New("invalid base64 text")
	}

	return string(decoded), nil
}

// IsValidBase64 checks if a string is valid base64
func (e *Encoder) IsValidBase64(text string) bool {
	if text == "" {
		return false
	}

	_, err := base64.StdEncoding.DecodeString(text)
	return err == nil
}

// EncodeBytes converts byte slice to base64
func (e *Encoder) EncodeBytes(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("data cannot be empty")
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

// DecodeBytes converts base64 text to byte slice
func (e *Encoder) DecodeBytes(encodedText string) ([]byte, error) {
	if encodedText == "" {
		return nil, errors.New("encoded text cannot be empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return nil, errors.New("invalid base64 text")
	}

	return decoded, nil
}
