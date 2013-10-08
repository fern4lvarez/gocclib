package main

import (
	"fmt"

	cc "github.com/fern4lvarez/gocclib"
)

func main() {
	api := cc.NewAPI()
	err := api.CreateToken("", "")
	if err != nil {
		panic(err)
	}
	app, err := api.ReadApp("fago")
	if err != nil {
		panic(err)
	}
	fmt.Println("App", app)
}
