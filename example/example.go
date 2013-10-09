package main

import (
	"fmt"
	"os"

	cc "github.com/fern4lvarez/gocclib"
)

func main() {
	api := cc.NewAPI()
	err := api.CreateToken("fa@cloudcontrol.de", "")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(0)
	}

	data, err := api.DeleteDeployment("fagocclib2", "fagocclib2")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(0)
	}
	fmt.Println(data)
}
