package main

import (
	"fmt"

	argModule "github.com/keyopto/kscraper/internal/arg"
	"github.com/keyopto/kscraper/internal/types"
)

func main() {
	var arg types.ArgumentCommand
	err := argModule.ParseArgs(&arg)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get all the elements of this page

	// search for the https addresses

	// check if this address has an error

	// Print results
	fmt.Println("hello")
}
