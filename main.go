package main

import (
	"fmt"
	"os"

	"github.com/aktsk/bqnotify/lib/runner"
)

var version = "0.3.0"

func main() {
	err := runner.Run("config.yaml")
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
