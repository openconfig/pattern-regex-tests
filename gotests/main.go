package main

import (
	"flag"
	"log"

	"github.com/openconfig/pattern-regex-tests/gotests/patterncheck"
)

var (
	modelRoot string
	verbose   bool
)

func init() {
	flag.StringVar(&modelRoot, "model-root", "", "OpenConfig models root directory")
	flag.BoolVar(&verbose, "verbose", false, "print extra messages")
}

func main() {
	flag.Parse()

	if modelRoot == "" {
		log.Fatalf("Must supply model-root path")
	}

	patterncheck.Verbose = verbose

	if err := patterncheck.CheckRegexps(flag.Args(), []string{modelRoot}); err != nil {
		log.Fatal(err)
	}
}
