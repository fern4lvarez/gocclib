/*
gocclib is a library for accessing the cloudControl API using Go.
Please read the API documentation: https://api.cloudcontrol.com/doc/

Basic usage example:

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

*/
package cclib

import (
	"encoding/json"
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
//
// * Application name
//
// * Application type (php, python, ruby, java, nodejs, custom)
//
// * Repository type (git, bzr)
//
// * Buildpack URL (for custom application type)
//
// Returns an interface with just created application details
// and an error if request does not success.
func (api *API) CreateApp(appName, appType, repositoryType, buildpackURL string) (interface{}, error) {
	app := url.Values{}
	app.Add("name", appName)
	app.Add("type", appType)
	app.Add("repository_type", repositoryType)

	if buildpackURL != "" {
		app.Add("buildpack_url", buildpackURL)
	}

	return api.Post("/app/", app)
}

// ReadApps reads applications of current user.
//
// Returns an interface object with applications details
// and an error if request does not success.
func (api *API) ReadApps() (interface{}, error) {
	return api.Get("/app/")
}

// ReadApp reads an application having:
//
// * Application name
//
// Returns an interface with application details
// and an error if request does not success.
func (api *API) ReadApp(appName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/", appName))
}

// DeleteApp deletes an application having:
//
// * Application name
//
// Returns an error if request does not success.
func (api *API) DeleteApp(appName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/", appName))
}

/*
	Deployments
*/

// CreateDeployment creates a deployment having:
//
// * Application name
//
// * Deployment name
//
// * Stack name, optional
//
// Returns an interface with just created deployment details
// and an error if request does not success.
func (api *API) CreateDeployment(appName, depName, stack string) (interface{}, error) {
	dep := url.Values{}
	if depName != "" {
		dep.Add("name", depName)
	}

	if stack != "" {
		dep.Add("stack", stack)
	}

	return api.Post(fmt.Sprintf("/app/%s/deployment/", appName), dep)
}

// ReadDeployment reads a deployment having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with deployment details
// and an error if request does not success.
func (api *API) ReadDeployment(appName, depName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

// ReadDeployment reads deployment's users having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with deployment's users details
// and an error if request does not success.
func (api *API) ReadDeploymentUsers(appName, depName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/user/", appName, depName))
}

// UpdateDeployment updates deployment having:
// * Application name
//
// * Deployment name
//
// * Version to pull from the branch. Defaults to the last version if blank.
//
// * Billing account name of one of the users
//
// * Stack name
//
// * Number of containers constantly spawned: from 1 to 8
//
// * Size of containers: from 1 (128MB) to 8 (1024MB)
//
// Returns an interface with the updated deployment details
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

	return api.Put(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName), dep)
}

// DeleteDeployment deletes a deployment having:
//
// * Application name
//
// * Deployment name
//
// Returns an error if request does not success.
func (api *API) DeleteDeployment(appName, depName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
}

/*
	Aliases
*/

// CreateAlias creates an alias having:
//
// * Application name
//
// * Alias name
//
// * Deployment name
//
// Returns an interface with just created alias details
// and an error if request does not success.
func (api *API) CreateAlias(appName, aliasName, depName string) (interface{}, error) {
	alias := url.Values{}
	alias.Add("name", aliasName)

	return api.Post(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName), alias)
}

// ReadAliases reads all deployment's aliases having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with aliases details
// and an error if request does not success.
func (api *API) ReadAliases(appName, depName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName))
}

// ReadAlias reads a deployment's alias having:
//
// * Application name
//
// * Alias name
//
// * Deployment name
//
// Returns an interface with alias details
// and an error if request does not success.
func (api *API) ReadAlias(appName, aliasName, depName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
}

// DeleteAlias removes an alias from a deployment having:
//
// * Application name
//
// * Alias name
//
// * Deployment name
//
// Returns an error if request does not success.
func (api *API) DeleteAlias(appName, aliasName, depName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
}

/*
	Workers
*/

// CreateWorker executes a worker to a deployment having:
//
// * Application name
//
// * Deployment name
//
// * Worker command
//
// * Parameters, optional
//
// * Size, optional
//
// Returns an interface with just executed worker details
// and an error if request does not success.
func (api *API) CreateWorker(appName, depName, command, params, size string) (interface{}, error) {
	worker := url.Values{}
	worker.Add("command", command)

	if params != "" {
		worker.Add("params", params)
	}

	if size != "" {
		worker.Add("size", size)
	}

	return api.Post(fmt.Sprintf("/app/%s/deployment/%s/worker/", appName, depName), worker)
}

// ReadWorkers reads all deployment's workers having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with workers details
// and an error if request does not success.
func (api *API) ReadWorkers(appName, depName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/worker/", appName, depName))
}

// ReadWorker reads a deployment's worker having:
//
// * Application name
//
// * Deployment name
//
// * Worker ID
//
// Returns an interface with worker details
// and an error if request does not success.
func (api *API) ReadWorker(appName, depName, workerId string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/worker/%s/", appName, depName, workerId))
}

// DeleteWorker removes a worker from a deployment having:
//
// * Application name
//
// * Deployment name
//
// * Worker ID
//
// Returns an error if request does not success.
func (api *API) DeleteWorker(appName, depName, workerId string) error {
	return api.Delete(fmt.Sprintf("/app/%s/deployment/%s/worker/%s/", appName, depName, workerId))
}

/*
	Cronjobs
*/

// CreateCronjob adds a cronjob to a deployment having:
//
// * Application name
//
// * Deployment name
//
// * Cronjob URL
//
// Returns an interface with just created cronjob details
// and an error if request does not success.
func (api *API) CreateCronjob(appName, depName, urlJob string) (interface{}, error) {
	cronjob := url.Values{}
	cronjob.Add("url", urlJob)

	return api.Post(fmt.Sprintf("/app/%s/deployment/%s/cron/", appName, depName), cronjob)
}

// ReadCronjobs reads all deployment's cronjobs having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with cronjobs details
// and an error if request does not success.
func (api *API) ReadCronjobs(appName, depName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/cron/", appName, depName))
}

// ReadCronjob reads a deployment's cronjob having:
//
// * Application name
//
// * Deployment name
//
// * Cronjob ID
//
// Returns an interface with worker details
// and an error if request does not success.
func (api *API) ReadCronjob(appName, depName, cronjobId string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/cron/%s/", appName, depName, cronjobId))
}

// DeleteCronjob removes a cronjob from a deployment having:
//
// * Application name
//
// * Deployment name
//
// * Cronjob ID
//
// Returns an error if request does not success.
func (api *API) DeleteCronjob(appName, depName, cronjobId string) error {
	return api.Delete(fmt.Sprintf("/app/%s/deployment/%s/cron/%s/", appName, depName, cronjobId))
}

/*
	Addons
*/

// CreateAddon creates an addon having:
//
// * Application name
//
// * Deployment name
//
// * Addon name
//
// * Options as a pointer to a map of string to strings, optional
//
// Returns an interface with just created addon details
// and an error if request does not success.
func (api *API) CreateAddon(appName, depName, addonName string, options *map[string]string) (interface{}, error) {
	addon := url.Values{}
	addon.Add("addon", addonName)

	o, err := json.Marshal(&options)
	if err != nil {
		return nil, err
	}

	addon.Add("options", string(o))

	return api.Post(fmt.Sprintf("/app/%s/deployment/%s/addon/", appName, depName), addon)
}

// ReadCronjobs reads all deployment's addons having:
//
// * Application name
//
// * Deployment name
//
// If Application and Deployment names are empty,
// it returns of available addons.
// Otherwise it returns an interface with deployment's addons details
// and an error if request does not success.
func (api *API) ReadAddons(appName, depName string) (interface{}, error) {
	if appName != "" && depName != "" {
		return api.Get(fmt.Sprintf("/app/%s/deployment/%s/addon/", appName, depName))
	}
	return api.Get("/addon/")
}

// ReadAddon reads a deployment's addon having:
//
// * Application name
//
// * Deployment name
//
// * Addon name
//
// Returns an interface with addon details
// and an error if request does not success.
func (api *API) ReadAddon(appName, depName, addonName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/deployment/%s/addon/%s/", appName, depName, addonName))
}

// UpdateAddon updates addon having:
//
// * Application name
//
// * Deployment name
//
// * Current Addon name
//
// * New Addon name to update to
//
// Returns an interface with the updated addon details
// and an error if request does not success.
func (api *API) UpdateAddon(appName, depName, addonName, addonNameToUpdateTo string) (interface{}, error) {
	if depName == "" {
		depName = "default"
	}

	addon := url.Values{}
	addon.Add("addon", addonNameToUpdateTo)

	return api.Put(fmt.Sprintf("/app/%s/deployment/%s/addon/%s/", appName, depName, addonName), addon)
}

// DeleteAddon deletes an addon having:
//
// * Application name
//
// * Deployment name
//
// * Addon name
//
// Returns an error if request does not success.
func (api *API) DeleteAddon(appName, depName, addonName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/deployment/%s/addon/%s/", appName, depName, addonName))
}

/*
	App Users
*/

// CreateAppUser adds an user to an application having:
//
// * Application name
//
// * User email
//
// * User role, optional
//
// Returns an interface with just created user details
// and an error if request does not success.
func (api *API) CreateAppUser(appName, userEmail, role string) (interface{}, error) {
	user := url.Values{}
	user.Add("email", userEmail)

	if role != "" {
		user.Add("role", role)
	}

	return api.Post(fmt.Sprintf("/app/%s/user/", appName), user)
}

// ReadAppUsers reads all application's users having:
//
// * Application name
//
// Returns an interface with users details
// and an error if request does not success.
func (api *API) ReadAppUsers(appName string) (interface{}, error) {
	return api.Get(fmt.Sprintf("/app/%s/user/", appName))
}

// DeleteAppUser removes an user from an application having:
//
// * Application name
//
// * User name
//
// Returns an error if request does not success.
func (api *API) DeleteAppUser(appName, userName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/user/%s/", appName, userName))
}

/*
	Deployment Users
*/

// CreateDeploymentUser adds an user to a deployment having:
//
// * Application name
//
// * Deployment name
//
// * User email
//
// * User role, optional
//
// Returns an interface with just created user details
// and an error if request does not success.
func (api *API) CreateDeploymentUser(appName, depName, userEmail, role string) (interface{}, error) {
	user := url.Values{}
	user.Add("email", userEmail)

	if role != "" {
		user.Add("role", role)
	}

	return api.Post(fmt.Sprintf("/app/%s/deployment/%s/user/", appName, depName), user)
}

// DeleteDeploymentUser removes an user from a deployment having:
//
// * Application name
//
// * Deployment name
//
// * User name
//
// Returns an error if request does not success.
func (api *API) DeleteDeploymentUser(appName, depName, userName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/deployment/%s/user/%s/", appName, depName, userName))
}

/*
	Request wrappers
*/

// Get makes a GET request having a resource and data.
//
// Returns an interface with the requested object
// and an error if request does not success.
func (api *API) Get(resource string) (interface{}, error) {
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

// Post makes a POST request having a resource and data.
//
// Returns an interface with the new object
// and an error if request does not success.
func (api *API) Post(resource string, data url.Values) (interface{}, error) {
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

// Put makes a PUT request having a resource and data.
//
// Returns an interface with the updated object
// and an error if request does not success.
func (api *API) Put(resource string, data url.Values) (interface{}, error) {
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

// Delete makes a DELETE request having a resource.
//
// Returns an error if request does not success.
func (api *API) Delete(resource string) error {
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
