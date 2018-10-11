package handler

import (
	"testing"
)

func TestGetExampleHash(t *testing.T) {
	h := NewHasher()
	str := "angryMonkey"
	result0 := h.Hash(str)
	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if result0 != expected {
		t.Errorf("Expected (%s) got (%s)", expected, result0)
	}
}

func TestGetSameHash(t *testing.T) {
	h := NewHasher()
	str := "angryMonkey"
	result0 := h.Hash(str)
	result1 := h.Hash(str)
	if result0 != result1 {
		t.Errorf("I don't know how hashes work")
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
