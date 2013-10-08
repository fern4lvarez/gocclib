package cclib

type User struct {
	username string
	email    string
	role     string
}

type Users []User

type Owner struct {
	username                  string
	is_active                 bool
	lastName                  string
	email                     string
	firstName                 string
	billingWithoutPaymentInfo bool
}

type BuildpackURL string

type DateModified string
type DateCreated string

type Repository string

type Name string

type Type struct {
	name Name
}

type Invitation struct {
	email       string
	dateCreated DateCreated
}

type Invitations []Invitation

type Deployment struct {
	name  string
	depID string
}

type Deployments []Deployment

type App struct {
	buildpackURL BuildpackURL
	dateModified DateModified
	repository   Repository
	owner        Owner
	name         Name
	ty           Type
	invitations  Invitations
	deployments  Deployments
	dateCreated  DateCreated
	users        Users
}

type Apps []App
