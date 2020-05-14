package main

import (
	"fmt"
	"os"

	"github.com/masu-mi/goone/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
