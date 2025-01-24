package kevlar

import (
	"crypto/sha256"
	"io"
)

func Sha256(reader io.Reader) ([]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, reader)
	return h.Sum(nil), err
}
