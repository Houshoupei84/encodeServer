package handler

import (
	"testing"
)

func TestGetHash(t *testing.T) {
	h := NewHasher()
	str := "idc"
	result0 := h.Hash(str)
	result1 := h.Hash(str)
	if len(result0) != len(result1) {
		t.Errorf("I don't know how hashes work")
	}

	for i := range result0 {
		if result0[i] != result1[i] {
			t.Errorf("I'm confused")
		}
	}
}

func TestGetOtherHash(t *testing.T) {
	h := NewHasher()
	str := "idc if you know how long my hash is"
	result0 := h.Hash(str)
	result1 := h.Hash(str)
	if len(result0) != len(result1) {
		t.Errorf("I don't know how hashes work")
	}

	for i := range result0 {
		if result0[i] != result1[i] {
			t.Errorf("I'm confused")
		}
	}
}
