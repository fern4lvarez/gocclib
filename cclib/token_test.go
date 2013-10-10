package cclib

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	// Given
	token := Token{}
	b := []byte(`{"token":"abcdefghijklmnopqrstuvxyz"}`)

	// When
	err := token.Decode(b)

	// Then
	if err != nil {
		t.Errorf(msgFail, "Decode", nil, err)
	}
	if token.Key() != "abcdefghijklmnopqrstuvxyz" {
		t.Errorf(msgFail, "Decode", token, Token{"token": token.Key()})
	}
}

func TestEncode(t *testing.T) {
	// Given
	token := Token{"token": "abcdefghijklmnopqrstuvxyz"}

	// When
	b, err := token.Encode()

	// Then
	if err != nil {
		t.Errorf(msgFail, "Encode", nil, err)
	}
	if encoded := []byte(`{"token":"abcdefghijklmnopqrstuvxyz"}`); bytes.Compare(b, encoded) != 0 {
		t.Errorf(msgFail, "Encode", b, encoded)
	}
}

func TestKey(t *testing.T) {
	// Given
	token := Token{"token": "abcdefghijklmnopqrstuvxyz"}

	// When
	key := token.Key()

	// Then
	if key != "abcdefghijklmnopqrstuvxyz" {
		t.Errorf(msgFail, "Key", key, "abcdefghijklmnopqrstuvxyz")
	}
}
