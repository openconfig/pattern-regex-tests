// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Binary gotests is a commandline utility for running tests on the
// oc-ext:posix-pattern statements of a YANG file. See testdata for test
// extensions and how it's used.
package main

import (
	"flag"

	log "github.com/golang/glog"
	"github.com/openconfig/pattern-regex-tests/gotests/patterncheck"
	"github.com/openconfig/ygot/util"
)

var (
	modelRoot = flag.String("model-root", "", "OpenConfig models root directory")
)

func main() {
	flag.Parse()

	if *modelRoot == "" {
		log.Error("Must supply model-root path")
	}

	if err := patterncheck.CheckRegexps(flag.Args(), []string{*modelRoot}); err != nil {
		if errors, ok := err.(util.Errors); ok {
			for _, err := range errors {
				log.Errorln(err)
			}
		} else {
			log.Exit(err)
		}
	}
}
