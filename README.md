# gocclib

[Documentation online](http://godoc.org/github.com/fern4lvarez/gocclib/cclib)

**cclib** is a wrapper for the cloudControl API written in Go.
Please read the API documentation: https://api.cloudcontrol.com/doc/

## Install (with GOPATH set on your machine)
----------

* Step 1: Get the `cclib` package

```
go get github.com/fern4lvarez/gocclib/cclib
```

* Step 2 (Optional): Run tests

```
$ cd $GOPATH/src/github.com/fern4lvarez/gocclib/cclib
$ go test -v ./...
```

##Usage

```go
package main

import (
  "fmt"
  "os"

  cc "github.com/fern4lvarez/gocclib/cclib"
)

func main() {
  api := cc.NewAPI()
  err := api.CreateTokenFromFile("path_to_creds_file")
  if err != nil {
    fmt.Printf("%s\n", err.Error())
    os.Exit(0)
  }

  apps := api.ReadApps()
  ...
}
```


##License
----------
gocclib is Apache licensed.