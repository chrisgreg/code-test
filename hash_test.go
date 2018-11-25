package main

import (
	"testing"
)

func TestHashFunction(t *testing.T) {
	t.Log("Hashing test")
	hash := generateHash("test")
	expectedHash := "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"
	if hash != expectedHash {
		t.Errorf("Expected hash of %s, but it was %s instead.", expectedHash, hash)
	}
}
