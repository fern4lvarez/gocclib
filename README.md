gocclib [![Build Status](https://travis-ci.org/fern4lvarez/gocclib.png)](https://travis-ci.org/fern4lvarez/gocclib)[![GoDoc](http://godoc.org/github.com/fern4lvarez/gocclib/cclib?status.png)](http://godoc.org/github.com/fern4lvarez/gocclib/cclib)
========

**cclib** is a wrapper for the cloudControl API written in Go.
Please read the API documentation: https://api.cloudcontrol.com/doc/

Install
-------

* Get the `cclib` package

```
go get -u github.com/fern4lvarez/gocclib/cclib
```

* Run tests

```
$ cd $GOPATH/src/github.com/fern4lvarez/gocclib/cclib
$ go test -v ./...
```

Usage
-----

~~~go
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
~~~

##License
----------
gocclib is Apache licensed.
