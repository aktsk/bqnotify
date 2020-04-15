package main

import (
	"fmt"
	"os"

	"github.com/aktsk/bqnotify/lib/runner"
)

var version = "0.2.0"

func main() {
	err := runner.Run()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
