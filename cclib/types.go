package cclib

type ApplicationType struct {
	Name string `mapstructure:"name"`
}

type Owner struct {
	Username  string `mapstructure:"username"`
	FirstName string `mapstructure:"first_name"`
	LastName  string `mapstructure:"last_name"`
	Email     string `mapstructure:"email"`
	IsActive  bool   `mapstructure:"is_active"`
}

type User struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Role     string `mapstructure:"role"`
}

type Stack struct {
	Name string `mapstructure:"name"`
}

type BilledAddon struct {
	Name  string  `mapstructure:"addon"`
	Hours int     `mapstructure:"hours"`
	Costs int     `mapstructure:"costs"`
	Until float64 `mapstructure:"until"`
}

type Boxes struct {
	Boxes     int     `mapstructure:"boxes"`
	Costs     float32 `mapstructure:"costs"`
	FreeBoxes int     `mapstructure:"free_boxes"`
	Until     float64 `mapstructure:"until"`
}

type SupportPlan struct {
	Name                  string `mapstructure:"name"`
	ThirtyDaysPrice       string `mapstructure:"thirty_days_price"`
	PriceInBillPercentage string `mapstructure:"price_in_bill_percentage"`
}

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

type Deployment struct {
	Name             string         `mapstructure:"name"`
	DepId            string         `mapstructure:"dep_id"`
	DefaultSubdomain string         `mapstructure:"default_subdomain"`
	Users            []User         `mapstructure:"users"`
	Stack            Stack          `mapstructure:"stack"`
	BilledAddons     []BilledAddon  `mapstructure:"billed_addons""`
	Version          string         `mapstructure:"version"`
	IsDefault        bool           `mapstructure:"is_default"`
	BilledBoxes      Boxes          `mapstructure:"boxes"`
	BillingAccount   BillingAccount `mapstructure:"billing_account"`
	State            string         `mapstructure:"state"`
	Containers       int            `mapstructure:"min_boxes"`
	Size             int            `mapstructure:"max_boxes"`
}

type Application struct {
	Name         string          `mapstructure:"name"`
	Type         ApplicationType `mapstructure:"type"`
	Owner        Owner           `mapstructure:"owner"`
	BuildpackUrl string          `mapstructure:"buildpack_url"`
	Users        []User          `mapstructure:"users"`
	Deployments  []Deployment    `mapstructure:"deployments"`
}

type Alias struct {
	Name               string `mapstructure:"name"`
	VerificationCode   string `mapstructure:"verification_code"`
	VerificationErrors int    `mapstructure:"verification_errors"`
	IsDefault          bool   `mapstructure:"is_default"`
	IsVerified         bool   `mapstructure:"is_verified"`
}

type Worker struct {
	Id      string `mapstructure:"wrk_id"`
	Command string `mapstructure:"command"`
}

type Cronjob struct {
	Id string `mapstructure:"job_id"`
}

type AddonOption struct {
	Name string `mapstructure:"name"`
}

type Settings map[string]interface{}

type ConfigVars map[string]string

type Addon struct {
	Name     string      `mapstructure:"name"`
	Option   AddonOption `mapstructure:"addon_option"`
	Settings Settings    `mapstructure:"settings"`
}

type Key struct {
	Id string `mapstructure:"key_id"`
}

type Log struct {
	Type    string
	Message string
	Time    float64
}
