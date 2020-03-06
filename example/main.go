package main

import (
	"fmt"
	"os"

	"github.com/mizzy/bqnotify/lib/runner"
)

func main() {
	err := runner.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
