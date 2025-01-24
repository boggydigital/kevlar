package kevlar_legacy

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func Sha256(reader io.Reader) (string, error) {
	h := sha256.New()
	var err error
	if _, err = io.Copy(h, reader); err == nil {
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	}
	return "", err
}
