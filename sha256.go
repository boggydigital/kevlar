package kvas

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// Sha256 computes SHA-256 hash of a bytes slice
func Sha256(value io.Reader) (string, error) {
	h := sha256.New()
	_, err := io.Copy(h, value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
