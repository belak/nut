package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var (
	codec = pflag.String("codec", "json", "Which codec to use for reading and writing. Possible values: json, raw")
)

func main() {
	pflag.Parse()

	args := pflag.Args()
	if len(args) < 1 {
		fmt.Println("Not enough arguments. Please provide the path to a db.")
		pflag.Usage()
		os.Exit(1)
	}

	repl, err := NewNutDBRepl(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = repl.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
