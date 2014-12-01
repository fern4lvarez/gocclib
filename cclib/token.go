package cclib

import (
	"encoding/json"
	"io/ioutil"
)

// Token is the generated security and temporal token
// which contains the Key and the date when it Expires
type Token struct {
	Key     string `json:"token"`
	Expires string `json:"expires"`
}

// NewToken returns a pointer to a Token
// given a string token and a expiring time
func NewToken(token string, expires string) *Token {
	return &Token{
		Key:     token,
		Expires: expires,
	}
}

// Decode decodes bytes into a token
func (token *Token) Decode(b []byte) (err error) {
	return json.Unmarshal(b, &token)
}

// Encode encodes a token into bytes
func (token Token) Encode() (b []byte, err error) {
	b, err = json.Marshal(token)
	return
}

// Write writes the Token in a file
// in json format given the file path
func (token *Token) Write(path string) error {
	b, err := token.Encode()
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, []byte(b), 0644); err != nil {
		return err
	}

	return nil
}

// Read reads a Token from a given path
func (token *Token) Read(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = token.Decode(b); err != nil {
		return err
	}

	return nil
}
