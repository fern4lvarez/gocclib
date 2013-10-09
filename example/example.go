package main

import (
	"fmt"
	"os"

	cc "github.com/fern4lvarez/gocclib"
)

func main() {
	api := cc.NewAPI()
	err := api.CreateToken("", "")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(0)
	}
}
