/*
gocclib is a library for accessing the cloudControl API using Go.
Please read the API documentation: https://api.cloudcontrol.com/doc/

Basic usage example:

```
package main

import (
	"os"

	cc "github.com/fern4lvarez/gocclib/cclib"
)

func main() {
	api := cc.NewAPI()
	err := api.CreateToken("user@email.org", "password")
	if err != nil {
		os.Exit(0)
	}
	...
}
```
*/
package cclib

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// An API is the entity that make calls and manage
// the cloudControl interface.
type API struct {
	cache string
	url   string
	token *Token
}

// NewAPI creates a new API instance.
func NewAPI() *API {
	return NewAPIToken(nil)
}

// NewAPIToken creates an API instance from a token.
func NewAPIToken(token *Token) *API {
	return &API{"", API_URL, token}
}

// RequiresToken returns an error if API has no token.
func (api API) RequiresToken() (e error) {
	if !api.CheckToken() {
		e = errors.New("Token required.")
	}
	return
}

// CreateTokenFromFile creates a token for an api from
// a credential file's path.
// This file must contain two lines, first for the password and
// second for the password. Returns an error if credentials file
// is not OK or if there is problems creating a token.
func (api *API) CreateTokenFromFile(filepath string) (err error) {
	email, password, err := readCredentialsFile(filepath)
	if err != nil {
		return err
	}

	return api.CreateToken(email, password)
}

// CreateToken creates a token for an api from
// a email and password.
// Returns an error if there is problems creating a token.
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

// CheckToken returns true if api has a token.
func (api API) CheckToken() bool {
	if api.Token() == nil {
		return false
	} else {
		return true
	}
}

// SetToken sets a token to an api.
func (api *API) SetToken(t *Token) {
	api.token = t
}

// Cache returns api's cache.
func (api API) Cache() string {
	return api.cache
}

// Url returns api's URL.
func (api API) Url() string {
	return api.url
}

// Token returns api's token.
func (api API) Token() *Token {
	return api.token
}

/*
	Applications
*/

// CreateApp creates an application having:
// * Application name
// * Application type (php, python, ruby, java, nodejs, custom)
// * Repository type (git, bzr)
// * Buildpack URL (for custom application type)
//
// Returns an interface object with just created application information
// and an error if request does not success.
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

// ReadApps reads applications of current user.
//
// Returns an interface object with applications information
// and an error if request does not success.
func (api *API) ReadApps() (interface{}, error) {
	return api.getRequest("/app/")
}

// ReadApp reads an application having:
// * Application name
//
// Returns an interface object with application information
// and an error if request does not success.
func (api *API) ReadApp(appName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/", appName))
}

// DeleteApp deletes an application having:
// * Application name
//
// Returns an error if request does not success.
func (api *API) DeleteApp(appName string) error {
	return api.deleteRequest(fmt.Sprintf("/app/%s/", appName))
}

/*
	Deployments
*/

// CreateDeployment creates a deployment having:
// * Application name
// * Deployment name
// * Stack name (Optional)
//
// Returns an interface object with just created deployment information
// and an error if request does not success.
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

// ReadDeployment reads a deployment having:
// * Application name
// * Deployment name
//
// Returns an interface object with deployment information
// and an error if request does not success.
func (api *API) ReadDeployment(appName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

// ReadDeployment reads deployment's users having:
// * Application name
// * Deployment name
//
// Returns an interface object with deployment's users information
// and an error if request does not success.
func (api *API) ReadDeploymentUsers(appName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/user/", appName, depName))
}

// UpdateDeployment updates deployment having:
// * Application name
// * Deployment name
// * Version to pull from the branch. Defaults to the last version if blank.
// * Billing account name of one of the users
// * Stack name
// * Number of containers constantly spawned: from 1 to 8
// * Size of containers: from 1 (128MB) to 8 (1024MB)
//
// Returns an interface object with the updated deployment information
// and an error if request does not success.
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

// DeleteDeployment deletes a deployment having:
// * Application name
// * Deployment name
//
// Returns an error if request does not success.
func (api *API) DeleteDeployment(appName, depName string) error {
	return api.deleteRequest(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

/*
	Aliases
*/

// CreateAlias creates an alias having:
// * Application name
// * Alias name
// * Deployment name
//
// Returns an interface object with just created alias information
// and an error if request does not success.
func (api *API) CreateAlias(appName, aliasName, depName string) (interface{}, error) {
	alias := url.Values{}
	alias.Add("name", aliasName)

	return api.postRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName), alias)
}

// ReadAliases reads all deployment's aliases having:
// * Application name
// * Deployment name
//
// Returns an interface object with aliases information
// and an error if request does not success.
func (api *API) ReadAliases(appName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName))
}

// ReadAlias reads a deployment's alias having:
// * Application name
// * Alias name
// * Deployment name
//
// Returns an interface object with aliase information
// and an error if request does not success.
func (api *API) ReadAlias(appName, aliasName, depName string) (interface{}, error) {
	return api.getRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
}

// DeleteAlias deletes a deployment's alias having:
// * Application name
// * Alias name
// * Deployment name
//
// Returns an error if request does not success.
func (api *API) DeleteAlias(appName, aliasName, depName string) error {
	return api.deleteRequest(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
}

/*
	Request wrappers
*/

// getRequests makes a GET request having a resource.
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

// postRequests makes a POST request having a resource and data.
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

// putRequests makes a PUT request having a resource and data.
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

// deleteRequests makes a DELETE request having a resource.
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
