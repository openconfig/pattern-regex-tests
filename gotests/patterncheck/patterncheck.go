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

// Package patterncheck provides functions to execute regular expression checks
// defined in the OpenConfig yang modules against static test data to ensure
// the yang regexes have their intended effect.
package patterncheck

import (
	"fmt"
	"regexp"

	log "github.com/golang/glog"

	"github.com/openconfig/goyang/pkg/yangentry"
	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/ygot/util"
)

var (
	passCaseExt = regexp.MustCompile(`\w+:pattern-test-pass`)
	failCaseExt = regexp.MustCompile(`\w+:pattern-test-fail`)
)

// CheckRegexps tests mock input data against a set of leaves that have pattern
// test cases specified for them. It ensures that the regexp compiles as a
// POSIX regular expression according to the OpenConfig style guide.
func CheckRegexps(yangfiles, paths []string) ([]string, error) {
	yangE, errs := yangentry.Parse(yangfiles, paths)
	if len(errs) != 0 {
		return nil, fmt.Errorf("could not parse modules: %v", errs)
	}
	if len(yangE) == 0 {
		return nil, fmt.Errorf("did not parse any modules")
	}

	var patternErrs util.Errors
	var allFailMessages []string
	for _, mod := range yangE {
		for _, entry := range mod.Dir {
			if failMessages, err := checkEntryPatterns(entry); err != nil {
				patternErrs = util.AppendErr(patternErrs, err)
			} else {
				allFailMessages = append(allFailMessages, failMessages...)
			}
		}
	}
	if len(patternErrs) != 0 {
		return nil, patternErrs
	}
	return allFailMessages, nil
}

func checkEntryPatterns(entry *yang.Entry) ([]string, error) {
	if entry.Kind != yang.LeafEntry {
		return nil, nil
	}

	if len(entry.Errors) != 0 {
		return nil, fmt.Errorf("entry had associated errors: %v", entry.Errors)
	}

	var failMessages []string
	for _, ext := range entry.Exts {
		var wantMatch bool
		switch {
		case passCaseExt.MatchString(ext.Keyword):
			wantMatch = true
		case failCaseExt.MatchString(ext.Keyword):
			wantMatch = false
		default:
			continue
		}

		var gotMatch bool
		if len(entry.Type.Type) == 0 {
			var err error
			if gotMatch, err = checkPatterns(ext.Argument, entry.Type.POSIXPattern); err != nil {
				return nil, err
			}
		} else {
			// Handle unions.
			// We don't want to short-circuit the match
			// calculation, as we want to make sure all
			// POSIX patterns compile.
			for _, membertype := range entry.Type.Type {
				if membertype.Kind != yang.Ystring {
					// Only do the test when the pattern is a string.
					continue
				}
				matches, err := checkPatterns(ext.Argument, membertype.POSIXPattern)
				if err != nil {
					return nil, err
				}
				gotMatch = gotMatch || matches
			}
		}

		matchDesc := fmt.Sprintf("| `%s` | `%s` | `%s` matched but shouldn't |", entry.Name, entry.Type.Name, ext.Argument)
		if !gotMatch {
			matchDesc = fmt.Sprintf("| `%s` | `%s` | `%s` did not match |", entry.Name, entry.Type.Name, ext.Argument)
		}

		if gotMatch != wantMatch {
			failMessages = append(failMessages, matchDesc)
		}
		log.Infof("pass: %s", matchDesc)
	}
	return failMessages, nil
}

// checkPatterns compiles all given POSIX patterns, and returns true if
// testData matches every one of them. If one or more patterns don't compile,
// it returns false, as well as an error containing all non-compiling patterns.
// https://tools.ietf.org/html/rfc7950#section-9.4.5:
// If the type has multiple "pattern" statements, the expressions are
// ANDed together, i.e., all such expressions have to match.
func checkPatterns(testData string, patterns []string) (bool, error) {
	var errs util.Errors
	matches := true
	for _, pattern := range patterns {
		if r, err := regexp.CompilePOSIX(pattern); err != nil {
			errs = util.AppendErr(errs, err)
			matches = false
		} else if !r.MatchString(testData) {
			// Don't return early -- want to make sure all POSIX
			// patterns compile.
			matches = false
		}
	}
	if len(errs) == 0 {
		return matches, nil
	}
	return matches, errs
}
