package cclib

// ApplicationType contains the type of the application:
type ApplicationType struct {
	// python, ruby, java, php, nodejs, and custom
	Name string `mapstructure:"name"`
}

// Owner contains information about an application owner
type Owner struct {
	Username  string `mapstructure:"username"`
	FirstName string `mapstructure:"first_name"`
	LastName  string `mapstructure:"last_name"`
	Email     string `mapstructure:"email"`
	IsActive  bool   `mapstructure:"is_active"`
}

// User contains information about an application user
type User struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	// owner, admin, and readonly
	Role string `mapstructure:"role"`
}

// Stack contains information about the stack version
type Stack struct {
	// luigi (lucid), pinky (precise)
	Name string `mapstructure:"name"`
}

// BilledAddon contains information about the billing
// of an add-on
type BilledAddon struct {
	Name  string  `mapstructure:"addon"`
	Hours int     `mapstructure:"hours"`
	Costs int     `mapstructure:"costs"`
	Until float64 `mapstructure:"until"`
}

// Boxes contains information about the billing of
// a deployment
type Boxes struct {
	Boxes     int     `mapstructure:"boxes"`
	Costs     float32 `mapstructure:"costs"`
	FreeBoxes int     `mapstructure:"free_boxes"`
	Until     float64 `mapstructure:"until"`
}

// SupportPlan contains information about a support plan
type SupportPlan struct {
	Name                  string `mapstructure:"name"`
	ThirtyDaysPrice       string `mapstructure:"thirty_days_price"`
	PriceInBillPercentage string `mapstructure:"price_in_bill_percentage"`
}

// BillingAccount contains information about a billing account
type BillingAccount struct {
	Default     bool        `mapstructure:"default"`
	Email       string      `mapstructure:"email"`
	PostalCode  string      `mapstructure:"postal_code"`
	Title       string      `mapstructure:"title"`
	Name        string      `mapstructure:"name"`
	FirstName   string      `mapstructure:"first_name"`
	SecondName  string      `mapstructure:"second_name"`
	User        User        `mapstructure:"user"`
	Company     string      `mapstructure:"company"`
	Country     string      `mapstructure:"country"`
	SupportPlan SupportPlan `mapstructure:"support_plan"`
}

// Deployment contains information about a deployment
type Deployment struct {
	Name string `mapstructure:"name",json:"name",yaml:"name"`
	// Id follows the format `depxxxxxxxx`
	Id               string         `mapstructure:"dep_id",json:"dep_id",yaml:"dep_id"`
	DefaultSubdomain string         `mapstructure:"default_subdomain",json:"default_subdomain,omitempty",yaml:"default_subdomain,omitempty"`
	Users            []User         `mapstructure:"users",json:"users,omitempty",yaml:"users,omitempty"`
	Stack            Stack          `mapstructure:"stack",json:"stack,omitempty",yaml:"stack,omitempty"`
	BilledAddons     []BilledAddon  `mapstructure:"billed_addons"",json:"billed_addons,omitempty",yaml:"billed_addons,omitempty"`
	Version          string         `mapstructure:"version",json:"version,omitempty",yaml:"version,omitempty"`
	IsDefault        bool           `mapstructure:"is_default",json:"is_default,omitempty",yaml:"is_default,omitempty"`
	BilledBoxes      Boxes          `mapstructure:"boxes",json:"boxes,omitempty",yaml:"boxes,omitempty"`
	BillingAccount   BillingAccount `mapstructure:"billing_account",json:"billing_account,omitempty",yaml:"billing_account,omitempty"`
	State            string         `mapstructure:"state",json:"state,omitempty",yaml:"state,omitempty"`
	// Containers mean the number of containers running per deployment
	Containers int `mapstructure:"min_boxes",json:"min_boxes,omitempty",yaml:"min_boxes,omitempty"`
	// Size of the container memory: 1->128MB, 2->256MB, ..., 8 -> 1024MB
	Size int `mapstructure:"max_boxes",json:"max_boxes,omitempty",yaml:"max_boxes,omitempty"`
}

// Application contains information about an application
type Application struct {
	Name  string          `mapstructure:"name",json:"name",yaml:"name"`
	Type  ApplicationType `mapstructure:"type",json:"type",yaml:"type"`
	Owner Owner           `mapstructure:"owner",json:"owner,omitempty",yaml:"owner,omitempty"`
	// BuildpackUrl is empty unless Type is `custom`
	BuildpackUrl string       `mapstructure:"buildpack_url",json:"buildpack_url,omitempty",yaml:"buildpack_url,omitempty"`
	Users        []User       `mapstructure:"users",json:"users,omitempty",yaml:"users,omitempty"`
	Deployments  []Deployment `mapstructure:"deployments",json:"deployments",yaml:"deployments"`
}

// Alias contains information about a deployment alias
type Alias struct {
	Name string `mapstructure:"name"`
	// VerificationCode is a code to be verified via TXT record
	VerificationCode string `mapstructure:"verification_code"`
	// VerificationErrors will be more than 0 if TXT record verification
	// failed
	VerificationErrors int `mapstructure:"verification_errors"`
	// IsDefault will be true if the alias is a native one:
	// * app_name.domain.com
	// * dep_name.domain.com, dep_name-domain.com
	IsDefault bool `mapstructure:"is_default"`
	// IsVerified is true if the TXT record verification succeeded
	IsVerified bool `mapstructure:"is_verified"`
}

// Worker contains information about a worker
type Worker struct {
	// Id follows the format `wrkxxxxxxxx`
	Id string `mapstructure:"wrk_id"`
	// Command contains the command the worker is executed with via Procfile
	Command string `mapstructure:"command"`
}

// Cronjob contains information about a cronjob
type Cronjob struct {
	// Id follows the format `jobxxxxxxxx`
	Id string `mapstructure:"job_id"`
}

// AddonOption contains information about an add-on option
type AddonOption struct {
	// Name follows the format ADDON_NAME.OPTION_NAME
	Name string `mapstructure:"name"`
}

// Setting contains the settings or options to create or
// update some of the add-on. They will be normalized to JSON
// by the API
type Settings map[string]interface{}

// Add contains the information about an add-on
type Addon struct {
	Name     string      `mapstructure:"name"`
	Option   AddonOption `mapstructure:"addon_option"`
	Settings Settings    `mapstructure:"settings"`
}

// Key contains the information about a user public key
type Key struct {
	// Id follows the format of a random string of 10 chars
	Id string `mapstructure:"key_id"`
}

// Log contains the information about a log entry
type Log struct {
	// error, deploy, and access
	Type    string
	Message string
	Time    float64
}
