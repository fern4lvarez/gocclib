package cclib

import (
	"encoding/json"
)

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
