package cclib

import (
	"encoding/json"
)

type Token map[string]string

func (token *Token) Decode(b []byte) (err error) {
	return json.Unmarshal(b, &token)
}

func (token Token) Encode() (b []byte, err error) {
	b, err = json.Marshal(token)
	return
}

func (token Token) Key() string {
	return token["token"]
}
