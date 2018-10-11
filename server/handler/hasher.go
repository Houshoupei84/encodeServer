package handler

import (
	"crypto/sha512"
	"encoding/base64"
)

type Hasher struct {
}

func NewHasher() *Hasher {
	h := new(Hasher)
	return h
}

func (h *Hasher) Hash(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(hasher.Sum([]byte(nil)))
}
