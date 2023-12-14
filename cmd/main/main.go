package main

import (
	"fmt"
	"os"

	"github.com/DEXPRO-Solutions-GmbH/xd/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
