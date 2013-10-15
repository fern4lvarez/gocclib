package main

import (
	"fmt"
	"os"

	cc "github.com/fern4lvarez/gocclib/cclib"
)

func main() {
	api := cc.NewAPI()
	err := api.CreateTokenFromFile("/home/fa/.cc")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(0)
	}

	data, err := api.ReadWorkers("faworkerapp", "default")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(0)
	}
	fmt.Println(data)
}
