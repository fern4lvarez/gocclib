/*
gocclib is a library for accessing the cloudControl API using Go.
Please read the API documentation: https://api.cloudcontrol.com/doc/

Basic usage example:

	package main

	import (
		"fmt"
		"os"

		cc "github.com/fern4lvarez/gocclib/cclib"
	)

	func main() {
		// Create API instance
		api := cc.NewAPI()

		// Create new user
		john, err := api.CreateUser("john", "john@example.org", "secret")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(john.Username, "was created.")

		// Activate User with the code provided by email
		john, err = api.ActivateUser("json", "activationcode")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(john.Username, "is active.")

		// Basic authentication to API
		err = api.CreateTokenFromFile("filepath")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Authenticated as", john.Username)

		// Create Application
		myapp, err := api.CreateApplication("myapp", "ruby", "git", "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(myapp.Name, "was created.")

		// Add user to the application
		_, err = api.CreateAppUser("myapp", "john@example.org", "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(john.Username, "was added to", myapp.Name)
	}

*/
package cclib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	ms "github.com/mitchellh/mapstructure"
)

// An API is the entity that make calls and manage
// the cloudControl interface.
type API struct {
	Cache string
	Url   string
	Token *Token
}

// NewAPI creates a new API instance.
func NewAPI() *API {
	return NewAPIToken(nil)
}

// NewAPIToken creates an API instance from a token.
func NewAPIToken(token *Token) *API {
	if api_url := os.Getenv("CCTRL_API_URL"); api_url != "" {
		API_URL = api_url
	}
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
//
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
// Returns an error if there is any problem creating the token.
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
	if api.Token == nil {
		return false
	}

	return true
}

// SetToken sets a token to an api.
func (api *API) SetToken(t *Token) {
	api.Token = t
}

/*
	Applications
*/

// CreateApplication creates an application having:
//
// * Application name
//
// * Application type (php, python, ruby, java, nodejs, custom)
//
// * Repository type (git, bzr)
//
// * Buildpack URL (for custom application type)
//
// Returns an Application
// and an error if request does not success.
func (api *API) CreateApplication(appName, appType, repositoryType, buildpackURL string) (*Application, error) {
	appValues := url.Values{}
	appValues.Add("name", appName)
	appValues.Add("type", appType)
	appValues.Add("repository_type", repositoryType)

	if buildpackURL != "" {
		appValues.Add("buildpack_url", buildpackURL)
	}

	data, err := api.Post("/app/", appValues)
	return api.decodeApplication(data, err)
}

// ReadApplications reads applications of current user.
//
// Returns a list of Applications
// and an error if request does not success.
func (api *API) ReadApplications() (*[]Application, error) {
	data, err := api.Get("/app/")
	return api.decodeApplications(data, err)
}

// ReadApplication reads an application having:
//
// * Application name
//
// Returns an Application and
// an error if request does not success.
func (api *API) ReadApplication(appName string) (*Application, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/", appName))
	return api.decodeApplication(data, err)
}

// DeleteApplication deletes an application having:
//
// * Application name
//
// Returns an error if request does not success.
func (api *API) DeleteApplication(appName string) error {
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
// Returns the just created Deployment
// and an error if request does not success.
func (api *API) CreateDeployment(appName, depName, stack string) (*Deployment, error) {
	dep := url.Values{}
	if depName != "" {
		dep.Add("name", depName)
	}

	if stack != "" {
		dep.Add("stack", stack)
	}

	data, err := api.Post(fmt.Sprintf("/app/%s/deployment/", appName), dep)
	return api.decodeDeployment(data, err)
}

// ReadDeployment reads a deployment having:
//
// * Application name
//
// * Deployment name
//
// Returns a Deployment and
// an error if request does not success.
func (api *API) ReadDeployment(appName, depName string) (*Deployment, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName))
	return api.decodeDeployment(data, err)
}

// ReadDeployments reads all user deployments having:
//
// * Application name
//
// * Deployment name
//
// Returns a Deployment and
// an error if request does not success.
func (api *API) ReadDeployments(appName string) (*[]Deployment, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/", appName))
	return api.decodeDeployments(data, err)
}

// UpdateDeployment updates deployment having:
//
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
// Returns the updated Deployment
// and an error if request does not success.
func (api *API) UpdateDeployment(appName, depName, version, billingAccount, stack string, containers, size int) (*Deployment, error) {
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

	data, err := api.Put(fmt.Sprintf("/app/%s/deployment/%s/", appName, depName), dep)
	return api.decodeDeployment(data, err)
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
// Returns the just created Alias
// and an error if request does not success.
func (api *API) CreateAlias(appName, aliasName, depName string) (*Alias, error) {

	aliasValues := url.Values{}
	aliasValues.Add("name", aliasName)

	data, err := api.Post(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName), aliasValues)
	return api.decodeAlias(data, err)
}

// ReadAliases reads all deployment's aliases having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with aliases details
// and an error if request does not success.
func (api *API) ReadAliases(appName, depName string) (*[]Alias, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/alias/", appName, depName))
	return api.decodeAliases(data, err)
}

// ReadAlias reads a deployment's alias having:
//
// * Application name
//
// * Alias name
//
// * Deployment name
//
// Returns an Alias
// and an error if request does not success.
func (api *API) ReadAlias(appName, aliasName, depName string) (*Alias, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/alias/%s/", appName, depName, aliasName))
	return api.decodeAlias(data, err)

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
// Returns the just created Worker
// and an error if request does not success.
func (api *API) CreateWorker(appName, depName, command, params, size string) (*Worker, error) {
	workerValues := url.Values{}
	workerValues.Add("command", command)

	if params != "" {
		workerValues.Add("params", params)
	}

	if size != "" {
		workerValues.Add("size", size)
	}

	data, err := api.Post(fmt.Sprintf("/app/%s/deployment/%s/worker/", appName, depName), workerValues)
	return api.decodeWorker(data, err)
}

// ReadWorkers reads all deployment's workers having:
//
// * Application name
//
// * Deployment name
//
// Returns a list of workers
// and an error if request does not success.
func (api *API) ReadWorkers(appName, depName string) (*[]Worker, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/worker/", appName, depName))
	return api.decodeWorkers(data, err)
}

// ReadWorker reads a deployment's worker having:
//
// * Application name
//
// * Deployment name
//
// * Worker ID
//
// Returns a Worker
// and an error if request does not success.
func (api *API) ReadWorker(appName, depName, workerId string) (*Worker, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/worker/%s/", appName, depName, workerId))
	return api.decodeWorker(data, err)
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
// Returns the just created Cronjob
// and an error if request does not success.
func (api *API) CreateCronjob(appName, depName, urlJob string) (*Cronjob, error) {
	cronjobValues := url.Values{}
	cronjobValues.Add("url", urlJob)

	data, err := api.Post(fmt.Sprintf("/app/%s/deployment/%s/cron/", appName, depName), cronjobValues)
	return api.decodeCronjob(data, err)
}

// ReadCronjobs reads all deployment's cronjobs having:
//
// * Application name
//
// * Deployment name
//
// Returns a Cronjob
// and an error if request does not success.
func (api *API) ReadCronjobs(appName, depName string) (*[]Cronjob, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/cron/", appName, depName))
	return api.decodeCronjobs(data, err)
}

// ReadCronjob reads a deployment's cronjob having:
//
// * Application name
//
// * Deployment name
//
// * Cronjob ID
//
// Returns a Cronjob
// and an error if request does not success.
func (api *API) ReadCronjob(appName, depName, cronjobId string) (*Cronjob, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/cron/%s/", appName, depName, cronjobId))
	return api.decodeCronjob(data, err)
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
// * Settings for the add-on, optional. eg:
//		settings := &Settings{
//			"foo":   "bar",
//			"hello": "world",
//			"mybool": "true",
//			"mynumber": "42"
//		}
//
//
// Returns the just created Addon
// and an error if request does not success.
func (api *API) CreateAddon(appName, depName, addonName string, settings *Settings) (*Addon, error) {
	addonValues := url.Values{}
	addonValues.Add("addon", addonName)

	o, err := json.Marshal(&settings)
	if err != nil {
		return nil, err
	}

	addonValues.Add("options", string(o))

	data, err := api.Post(fmt.Sprintf("/app/%s/deployment/%s/addon/", appName, depName), addonValues)
	return api.decodeAddon(data, err)
}

// ReadCronjobs reads all deployment's addons having:
//
// * Application name
//
// * Deployment name
//
// If Application and Deployment names are empty,
// it returns of available addons.
// Otherwise it returns deployment's Addons
// and an error if request does not success.
func (api *API) ReadAddons(appName, depName string) (*[]Addon, error) {
	var data interface{}
	var err error

	if appName != "" && depName != "" {
		data, err = api.Get(fmt.Sprintf("/app/%s/deployment/%s/addon/", appName, depName))
	} else {
		data, err = api.Get("/addon/")
	}

	return api.decodeAddons(data, err)
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
func (api *API) ReadAddon(appName, depName, addonName string) (*Addon, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/addon/%s/", appName, depName, addonName))
	return api.decodeAddon(data, err)
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
// * Settings for the add-on, optional. eg:
//		settings := &Settings{
//			"foo":   "bar",
//			"hello": "world",
//			"mybool": "true",
//			"mynumber": "42"
//		}
//
// * Force update in case of conflict
//
// Returns the updated Addon
// and an error if request does not success.
func (api *API) UpdateAddon(appName, depName, addonName, addonNameToUpdateTo string, settings *Settings, force bool) (*Addon, error) {
	if depName == "" {
		depName = "default"
	}

	addonValues := url.Values{}
	addonValues.Add("addon", addonNameToUpdateTo)

	if settings != nil {
		s, err := json.Marshal(&settings)
		if err != nil {
			return nil, err
		}

		addonValues.Add("settings", string(s))
	}

	if force {
		addonValues.Add("force", "true")
	}

	data, err := api.Put(fmt.Sprintf("/app/%s/deployment/%s/addon/%s/", appName, depName, addonName), addonValues)
	return api.decodeAddon(data, err)
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
// Returns the just added User
// and an error if request does not success.
func (api *API) CreateAppUser(appName, userEmail, role string) (*User, error) {
	userValues := url.Values{}
	userValues.Add("email", userEmail)

	if role != "" {
		userValues.Add("role", role)
	}

	data, err := api.Post(fmt.Sprintf("/app/%s/user/", appName), userValues)
	return api.decodeUser(data, err)
}

// ReadAppUsers reads all application's users having:
//
// * Application name
//
// Returns a list of application Users
// and an error if request does not success.
func (api *API) ReadAppUsers(appName string) (*[]User, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/user/", appName))
	return api.decodeUsers(data, err)
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
// Returns the just created User
// and an error if request does not success.
func (api *API) CreateDeploymentUser(appName, depName, userEmail, role string) (*User, error) {
	userValues := url.Values{}
	userValues.Add("email", userEmail)

	if role != "" {
		userValues.Add("role", role)
	}

	data, err := api.Post(fmt.Sprintf("/app/%s/deployment/%s/user/", appName, depName), userValues)
	return api.decodeUser(data, err)
}

// ReadDeploymentUsers reads deployment's users having:
//
// * Application name
//
// * Deployment name
//
// Returns an interface with deployment's users details
// and an error if request does not success.
func (api *API) ReadDeploymentUsers(appName, depName string) (*[]User, error) {
	data, err := api.Get(fmt.Sprintf("/app/%s/deployment/%s/user/", appName, depName))
	return api.decodeUsers(data, err)
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
	Users
*/

// CreateUser creates a new user having:
//
// * User name
//
// * User email
//
// * User password
//
// Returns the just created User
// and an error if request does not success.
func (api *API) CreateUser(userName, userEmail, password string) (*User, error) {
	userValues := url.Values{}
	userValues.Add("username", userName)
	userValues.Add("email", userEmail)
	userValues.Add("password", password)

	data, err := api.Post("/user/", userValues)
	return api.decodeUser(data, err)
}

// ReadUsers gets users. Usually just your own.
//
// Returns a list of Users
// and an error if request does not success.
func (api *API) ReadUsers() (*[]User, error) {
	data, err := api.Get("/user/")
	return api.decodeUsers(data, err)
}

// ReadUser reads a user having:
//
// * User name
//
// Returns a User
// and an error if request does not success.
func (api *API) ReadUser(userName string) (*User, error) {
	data, err := api.Get(fmt.Sprintf("/user/%s/", userName))
	return api.decodeUser(data, err)
}

// ActivateUser activates a new user having:
//
// * User name
//
// * Activation code for activation after registration
//
// Returns the activated User
// and an error if request does not success.
func (api *API) ActivateUser(userName, activationCode string) (*User, error) {
	userValues := url.Values{}
	if activationCode != "" {
		userValues.Add("activation_code", activationCode)
	}

	data, err := api.Put(fmt.Sprintf("/user/%s/", userName), userValues)
	return api.decodeUser(data, err)
}

// UpdateUser updates an existing user having:
//
// * User name
//
// * First name, optional
//
// * Last name, optional
//
// * Password, optional
//
// * Email, optional
//
// NOTE: At least one of the arguments must be provided.
//
// Returns the updated User
// and an error if request does not success.
func (api *API) UpdateUser(userName, firstName, lastName, password, email string) (*User, error) {
	userValues := url.Values{}
	if firstName != "" {
		userValues.Add("first_name", firstName)
	}
	if lastName != "" {
		userValues.Add("last_name", lastName)
	}
	if password != "" {
		userValues.Add("password", password)
	}
	if email != "" {
		userValues.Add("email", email)
	}

	data, err := api.Put(fmt.Sprintf("/user/%s/", userName), userValues)
	return api.decodeUser(data, err)

}

// DeleteUser deletes as user having:
//
// * User name
//
// Returns an error if request does not success.
func (api *API) DeleteUser(userName string) error {
	return api.Delete(fmt.Sprintf("/app/%s/", userName))
}

/*
	Keys
*/

// CreateUserKey creates a new user's key having:
//
// * User name
//
// * Public key
//
// Returns the just created Key
// and an error if request does not success.
func (api *API) CreateUserKey(userName, publicKey string) (*Key, error) {
	keyValues := url.Values{}
	keyValues.Add("key", publicKey)

	data, err := api.Post(fmt.Sprintf("/user/%s/key/", userName), keyValues)
	return api.decodeKey(data, err)
}

// ReadUserKeys gets all user keys having:
//
// * UserName
//
// Returns a list of Keys
// and an error if request does not success.
func (api *API) ReadUserKeys(userName string) (*[]Key, error) {
	data, err := api.Get(fmt.Sprintf("/user/%s/key/", userName))
	return api.decodeKeys(data, err)
}

// ReadUserKey reads a user's key having:
//
// * User name
//
// * Key Id
//
// Returns the Key
// and an error if request does not success.
func (api *API) ReadUserKey(userName, keyId string) (*Key, error) {
	data, err := api.Get(fmt.Sprintf("/user/%s/key/%s/", userName, keyId))
	return api.decodeKey(data, err)
}

// DeleteUserKey deletes a user's key having:
//
// * User name
//
// * Key ID
//
// Returns an error if request does not success.
func (api *API) DeleteUserKey(userName, keyID string) error {
	return api.Delete(fmt.Sprintf("/user/%s/key/%s/", userName, keyID))
}

/*
	Logs
*/

// ReadLog gets a deployment's log having:
//
// * Application name
//
// * Deployment name
//
// * Log type: worker, error, access
//
// * Last time from where to read on, optional. Make use of a time.Time struct pointer.
//
// Returns a list of Logs
// and an error if request does not success.
func (api *API) ReadLog(appName, depName, logType string, lastTime *time.Time) (*[]Log, error) {
	var resource string

	if lastTime == nil {
		resource = fmt.Sprintf("/app/%s/deployment/%s/log/%s/", appName, depName, logType)
	} else {
		resource = fmt.Sprintf("/app/%s/deployment/%s/log/%s/?timestamp=%s/", appName, depName, logType, buildTimestamp(lastTime))
	}

	data, err := api.Get(resource)
	return api.decodeLogs(data, err)
}

/*
	Billing Accounts
*/

// CreateBillingAccount creates a new billing account having:
//
// * User name
//
// * Billing name
//
// * Billing data in url.Values format
//
// Returns just created BillingAccount
// and an error if request does not success.
func (api *API) CreateBillingAccount(userName, billingName string, billingData url.Values) (*BillingAccount, error) {
	data, err := api.Post(fmt.Sprintf("/user/%s/billing/%s/", userName, billingName), billingData)
	return api.decodeBillingAccount(data, err)

}

// ReadBillingAccounts gets all billing accounts from a user having:
//
// * User name
//
// Returns a list of user's BillingAccounts
// and an error if request does not success.
func (api *API) ReadBillingAccounts(userName string) (*[]BillingAccount, error) {
	data, err := api.Get(fmt.Sprintf("/user/%s/billing/", userName))
	return api.decodeBillingAccounts(data, err)
}

// UpdateBillingAccount updates an existing user's billing account having:
//
// * User name
//
// * Billing name
//
// * Billing data in url.Values format
//
// Returns updated user's BillingAccount
// and an error if request does not success.
func (api *API) UpdateBillingAccount(userName, billingName string, billingData url.Values) (*BillingAccount, error) {
	data, err := api.Put(fmt.Sprintf("/user/%s/billing/%s/", userName, billingName), billingData)
	return api.decodeBillingAccount(data, err)
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

	request := NewRequest("", "", api.Token)

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

	request := NewRequest("", "", api.Token)

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

	request := NewRequest("", "", api.Token)

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

	request := NewRequest("", "", api.Token)

	content, err := request.Delete(resource)
	if err != nil {
		return err
	}

	_, err = decodeContent(content)
	return err
}

/*
	Type decoders
*/

func (api *API) decodeApplication(data interface{}, err error) (*Application, error) {
	if err != nil {
		return nil, err
	}

	var application Application
	if err = ms.Decode(data, &application); err != nil {
		return nil, err
	}

	return &application, nil
}

func (api *API) decodeApplications(data interface{}, err error) (*[]Application, error) {
	if err != nil {
		return nil, err
	}

	var applications []Application
	if err = ms.Decode(data, &applications); err != nil {
		return nil, err
	}

	return &applications, nil
}

func (api *API) decodeDeployment(data interface{}, err error) (*Deployment, error) {
	if err != nil {
		return nil, err
	}

	var deployment Deployment
	if err = ms.Decode(data, &deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
}

func (api *API) decodeDeployments(data interface{}, err error) (*[]Deployment, error) {
	if err != nil {
		return nil, err
	}

	var deployments []Deployment
	if err = ms.Decode(data, &deployments); err != nil {
		return nil, err
	}

	return &deployments, nil
}

func (api *API) decodeAlias(data interface{}, err error) (*Alias, error) {
	if err != nil {
		return nil, err
	}

	var alias Alias
	if err = ms.Decode(data, &alias); err != nil {
		return nil, err
	}

	return &alias, nil
}

func (api *API) decodeAliases(data interface{}, err error) (*[]Alias, error) {
	if err != nil {
		return nil, err
	}

	var aliases []Alias
	if err = ms.Decode(data, &aliases); err != nil {
		return nil, err
	}

	return &aliases, nil
}

func (api *API) decodeWorker(data interface{}, err error) (*Worker, error) {
	if err != nil {
		return nil, err
	}

	var worker Worker
	if err = ms.Decode(data, &worker); err != nil {
		return nil, err
	}

	return &worker, nil
}

func (api *API) decodeWorkers(data interface{}, err error) (*[]Worker, error) {
	if err != nil {
		return nil, err
	}

	var workers []Worker
	if err = ms.Decode(data, &workers); err != nil {
		return nil, err
	}

	return &workers, nil
}

func (api *API) decodeCronjob(data interface{}, err error) (*Cronjob, error) {
	if err != nil {
		return nil, err
	}

	var cronjob Cronjob
	if err = ms.Decode(data, &cronjob); err != nil {
		return nil, err
	}

	return &cronjob, nil
}

func (api *API) decodeCronjobs(data interface{}, err error) (*[]Cronjob, error) {
	if err != nil {
		return nil, err
	}

	var cronjobs []Cronjob
	if err = ms.Decode(data, &cronjobs); err != nil {
		return nil, err
	}

	return &cronjobs, nil
}

func (api *API) decodeAddon(data interface{}, err error) (*Addon, error) {
	if err != nil {
		return nil, err
	}

	var addon Addon
	if err = ms.Decode(data, &addon); err != nil {
		return nil, err
	}

	return &addon, nil
}

func (api *API) decodeAddons(data interface{}, err error) (*[]Addon, error) {
	if err != nil {
		return nil, err
	}

	var addons []Addon
	if err = ms.Decode(data, &addons); err != nil {
		return nil, err
	}

	return &addons, nil
}

func (api *API) decodeUser(data interface{}, err error) (*User, error) {
	if err != nil {
		return nil, err
	}

	var user User
	if err = ms.Decode(data, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (api *API) decodeUsers(data interface{}, err error) (*[]User, error) {
	if err != nil {
		return nil, err
	}

	var users []User
	if err = ms.Decode(data, &users); err != nil {
		return nil, err
	}

	return &users, nil
}

func (api *API) decodeKey(data interface{}, err error) (*Key, error) {
	if err != nil {
		return nil, err
	}

	var ey Key
	if err = ms.Decode(data, &ey); err != nil {
		return nil, err
	}

	return &ey, nil
}

func (api *API) decodeKeys(data interface{}, err error) (*[]Key, error) {
	if err != nil {
		return nil, err
	}

	var eys []Key
	if err = ms.Decode(data, &eys); err != nil {
		return nil, err
	}

	return &eys, nil
}

func (api *API) decodeLog(data interface{}, err error) (*Log, error) {
	if err != nil {
		return nil, err
	}

	var log Log
	if err = ms.Decode(data, &log); err != nil {
		return nil, err
	}

	return &log, nil
}

func (api *API) decodeLogs(data interface{}, err error) (*[]Log, error) {
	if err != nil {
		return nil, err
	}

	var logs []Log
	if err = ms.Decode(data, &logs); err != nil {
		return nil, err
	}

	return &logs, nil
}

func (api *API) decodeBillingAccount(data interface{}, err error) (*BillingAccount, error) {
	if err != nil {
		return nil, err
	}

	var billingAccount BillingAccount
	if err = ms.Decode(data, &billingAccount); err != nil {
		return nil, err
	}

	return &billingAccount, nil
}

func (api *API) decodeBillingAccounts(data interface{}, err error) (*[]BillingAccount, error) {
	if err != nil {
		return nil, err
	}

	var billingAccounts []BillingAccount
	if err = ms.Decode(data, &billingAccounts); err != nil {
		return nil, err
	}

	return &billingAccounts, nil
}
