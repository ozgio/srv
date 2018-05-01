package main

import (
	"fmt"
	"os"

	"github.com/ozgio/srv/cmd"
)

func main() {
	_, err := cmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
