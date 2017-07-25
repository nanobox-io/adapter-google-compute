package main

import (
	"os"
	"fmt"

	"github.com/nanobox-io/nanobox-provider-golang"
)

func main() {
	err := provider.Start(GoogleCompute{}, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}