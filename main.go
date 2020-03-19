package main

import (
	"fmt"
	"os"

	"github.com/aktsk/bqnotify/lib/runner"
)

func main() {
	err := runner.Run()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
