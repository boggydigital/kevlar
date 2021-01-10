package internal

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(value []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
