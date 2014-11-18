package cclib

import (
	"encoding/json"
)

// Tokenizer is an interface for
// the authentication Token
// provided by the cloudControl platform
type Tokenizer interface {
	Decode([]byte) error
	Encode() ([]byte, error)
	Key() string
}

// Token is the generated security and temporal token
type Token map[string]string

// Decode decodes bytes into a token
func (token *Token) Decode(b []byte) (err error) {
	return json.Unmarshal(b, &token)
}

// Encode encodes a token into bytes
func (token Token) Encode() (b []byte, err error) {
	b, err = json.Marshal(token)
	return
}

// Key returns a token's key
func (token Token) Key() string {
	return token["token"]
}
