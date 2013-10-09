package cclib

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
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

func (api *API) DeleteApp(appName string) error {
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

func (api *API) UpdateDeployment(appName, depName, version, billingAccount, stack string, containers, size int) (interface{}, error) {
	if depName == "" {
		depName = "default"
	}

	dep := url.Values{}
	if version != "" {
		dep.Add("version", version)
	}
	if billingAccount != "" {
		dep.Add("billing_account", billingAccount)
	}
	if stack != "" {
		dep.Add("stack", stack)
	}
	if containers > 0 {
		dep.Add("min_boxes", strconv.Itoa(containers))
	}
	if size > 0 {
		dep.Add("max_boxes", strconv.Itoa(size))
	}
	if stack != "" {
		dep.Add("stack", stack)
	}

	return api.putRequest(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName), dep)
}

func (api *API) DeleteDeployment(appName, depName string) error {
	return api.deleteRequest(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

/*
	Aliases
*/
func (api *API) CreateAlias(appName, aliasName, depName string) (interface{}, error) {
	alias := url.Values{}
	alias.Add("name", aliasName)

	return api.postRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName), alias)
}

func (api *API) ReadAliases(appName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName))
}

func (api *API) ReadAlias(appName, aliasName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
}

func (api *API) DeleteAlias(appName, aliasName, depName string) error {
	return api.deleteRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
}

/*
	Request wrappers
*/
func (api *API) getRequest(resource string) (interface{}, error) {
	if err := api.RequiresToken(); err != nil {
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
	if err := api.RequiresToken(); err != nil {
		return nil, err
	}

	request := NewRequest("", "", api.Token())

	content, err := request.Post(resource, data)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}

func (api *API) putRequest(resource string, data url.Values) (interface{}, error) {
	if err := api.RequiresToken(); err != nil {
		return nil, err
	}

	request := NewRequest("", "", api.Token())

	content, err := request.Put(resource, data)
	if err != nil {
		return nil, err
	}

	return decodeContent(content)
}

func (api *API) deleteRequest(resource string) error {
	if err := api.RequiresToken(); err != nil {
		return err
	}

	request := NewRequest("", "", api.Token())

	content, err := request.Delete(resource)
	if err != nil {
		return err
	}

	_, err = decodeContent(content)
	return err
}
