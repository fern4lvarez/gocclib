package cclib

import (
	"errors"
	"fmt"
	"net/url"
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
	if err = token.Decode(content); err != nil {
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

/*
	Applications
*/
func (api *API) CreateApp(appName, appType, repositoryType, buildpackURL string) (interface{}, error) {
	app := url.Values{}
	app.Add("name", appName)
	app.Add("type", appType)
	app.Add("repository_type", repositoryType)

	if buildpackURL != "" {
		app.Add("buildpack_url", buildpackURL)
	}

	return api.postRequest("/app/", app)
}

func (api *API) ReadApps() (interface{}, error) {
	return api.getRequest("/app/")
}

func (api *API) ReadApp(appName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/", appName))
}

func (api *API) DeleteApp(appName string) (interface{}, error) {
	return api.deleteRequest(fmt.Sprintf("/app/%s/", appName))
}

/*
	Deployments
*/
func (api *API) CreateDeployment(appName, depName, stack string) (interface{}, error) {
	dep := url.Values{}
	if depName != "" {
		dep.Add("name", depName)
	}

	if stack != "" {
		dep.Add("stack", stack)
	}

	return api.postRequest(fmt.Sprintf("/app/%s/deployment/", appName), dep)
}

func (api *API) ReadDeployment(appName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

func (api *API) ReadDeploymentUsers(appName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/user/", appName, depName))
}

func (api *API) DeleteDeployment(appName, depName string) (interface{}, error) {
	return api.deleteRequest(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

/*
	Request wrappers
*/
func (api *API) getRequest(resource string) (interface{}, error) {
	err := api.RequiresToken()
	if err != nil {
		return nil, err
	}

	request := NewRequest("", "", api.Token())

	content, err := request.Get(resource)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}

func (api *API) postRequest(resource string, data url.Values) (interface{}, error) {
	err := api.RequiresToken()
	if err != nil {
		return nil, err
	}

	request := NewRequest("", "", api.Token())

	content, err := request.Post(resource, data)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}

func (api *API) deleteRequest(resource string) (interface{}, error) {
	err := api.RequiresToken()
	if err != nil {
		return nil, err
	}

	request := NewRequest("", "", api.Token())

	content, err := request.Delete(resource)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}
