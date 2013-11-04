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

	/*	fago, _ := api.ReadDeployment("fago", "default")
		fmt.Println(fago)*/

	// faKey, _ := api.ReadUserKeys("fernandoalvarez")
	//fmt.Println(faKey)

	ba, _ := api.ReadBillingAccounts("fernandoalvarez")
	fmt.Println(ba)
}
