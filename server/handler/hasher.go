package handler

import (
	"crypto/sha512"
	"encoding/base64"
	"hash"
)

type Hasher struct {
	impl hash.Hash
}

func NewHasher() *Hasher {
	h := new(Hasher)
	h.impl = sha512.New()
	return h
}

func (h *Hasher) Hash(s string) string {
	return base64.RawStdEncoding.EncodeToString(h.impl.Sum([]byte(s)))
}
