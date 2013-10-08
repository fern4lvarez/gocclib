package cclib

import (
	"errors"
	"fmt"
)

type API struct {
	cache string
	url   string
	token *Token
}

func NewAPI() *API {
	return NewAPIToken(nil)
}

func NewAPIToken(t *Token) *API {
	return &API{"", "", t}
}

func (api API) RequiresToken() (e error) {
	if !api.CheckToken() {
		e = errors.New("Token required.")
	}
	return
}

func (api *API) CreateToken(email string, password string) (err error) {
	request := NewRequest(email, password, nil)
	content, err := request.Post("/token/", nil)
	if err != nil {
		return err
	}

	var token Token
	err = token.Decode(content)
	if err != nil {
		return err
	}

	api.SetToken(&token)
	return
}

func (api API) CheckToken() bool {
	if api.Token() == nil {
		return false
	} else {
		return true
	}
}

func (api *API) SetToken(t *Token) {
	api.token = t
}

func (api API) Cache() string {
	return api.cache
}

func (api API) Url() string {
	return api.url
}

func (api API) Token() *Token {
	return api.token
}

func (api *API) ReadApps() (interface{}, error) {
	err := api.RequiresToken()
	if err != nil {
		return nil, err
	}

	resource := "/app/"
	request := NewRequest("", "", api.Token())

	content, err := request.Get(resource)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}

func (api *API) ReadApp(name string) (interface{}, error) {
	err := api.RequiresToken()
	if err != nil {
		return nil, err
	}

	resource := fmt.Sprintf("/app/%s/", name)
	request := NewRequest("", "", api.Token())

	content, err := request.Get(resource)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}
