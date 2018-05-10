package main

import (
	"fmt"
	"os"

	"github.com/ozgio/srv/cmd"
)

var VERSION = "dev"

func main() {
	_, err := cmd.Execute(VERSION)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
