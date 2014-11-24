package cclib

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	// Given
	b := []byte(`{"token":"abcdefghijklmnopqrstuvxyz","expires":"2014-11-24T16:39:54.450"}`)
	expectedToken := Token{
		Key:     "abcdefghijklmnopqrstuvxyz",
		Expires: "2014-11-24T16:39:54.450",
	}

	// When
	var token Token
	err := token.Decode(b)

	// Then
	if err != nil {
		t.Errorf(msgFail, "Decode", nil, err)
	}
	if token.Key != "abcdefghijklmnopqrstuvxyz" {
		t.Errorf(msgFail, "Decode", token, expectedToken.Key)
	}
	if token.Expires != "2014-11-24T16:39:54.450" {
		t.Errorf(msgFail, "Decode", token, expectedToken.Expires)
	}
}

func TestEncode(t *testing.T) {
	// Given
	token := Token{
		Key:     "abcdefghijklmnopqrstuvxyz",
		Expires: "2014-11-24T16:39:54.450",
	}
	expectedEncoded := []byte(`{"token":"abcdefghijklmnopqrstuvxyz","expires":"2014-11-24T16:39:54.450"}`)

	// When
	b, err := token.Encode()

	// Then
	if err != nil {
		t.Errorf(msgFail, "Encode", nil, err)
	}
	if bytes.Compare(b, expectedEncoded) != 0 {
		t.Errorf(msgFail, "Encode", b, expectedEncoded)
	}
}

func TestWrite(t *testing.T) {
	// Given
	token := Token{
		Key:     "abcdefghijklmnopqrstuvxyz",
		Expires: "2014-11-24T16:39:54.450",
	}
	location := "token.json"
	defer os.Remove(location)

	var expectedToken Token

	// When
	err := token.Write(location)

	// Then
	if err != nil {
		t.Errorf(msgFail, "Write", nil, err)
	}

	data, err := ioutil.ReadFile(location)
	if err != nil {
		t.Errorf(msgFail, "Write", nil, err)
	}

	expectedToken.Decode(data)

	if !reflect.DeepEqual(token, expectedToken) {
		t.Errorf(msgFail, "Write", expectedToken, token)
	}
}

func TestRead(t *testing.T) {
	// Given
	var token Token
	b := []byte(`{"token":"abcdefghijklmnopqrstuvxyz","expires":"2014-11-24T16:39:54.450"}`)
	expectedToken := Token{
		Key:     "abcdefghijklmnopqrstuvxyz",
		Expires: "2014-11-24T16:39:54.450",
	}
	location := "token.json"
	ioutil.WriteFile(location, []byte(b), 0644)
	defer os.Remove(location)

	// When
	err := token.Read(location)

	// Then
	if err != nil {
		t.Errorf(msgFail, "Read", nil, err)
	}

	if !reflect.DeepEqual(token, expectedToken) {
		t.Errorf(msgFail, "Read", expectedToken, token)
	}
}
