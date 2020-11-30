package gotests

import (
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/openconfig/pattern-regex-tests/gotests/gooctest"
)

const (
	testModule = "regexp-test"
	testfile   = "../testdata/regexp-test.yang"
)

var (
	passCaseExt = regexp.MustCompile(`\w:pattern-test-pass`)
	failCaseExt = regexp.MustCompile(`\w:pattern-test-fail`)
)

var ocdir string

// YANGLeaf is a structure sued to describe a particular leaf of YANG schema.
type YANGLeaf struct {
	module string
	name   string
}

// RegexpTest specifies a test case for a particular regular expression check.
type RegexpTest struct {
	inData    string
	wantMatch bool
}

// TestRegexps tests mock input data against a set of leaves that have patterns
// specified for them. It ensures that the regexp compiles as a POSIX regular
// expression according to the OpenConfig style guide.
func TestRegexps(t *testing.T) {
	modules := []string{testfile}

	yangE, errs := gooctest.ProcessModules(modules, []string{ocdir})
	if len(errs) != 0 {
		t.Fatalf("could not parse modules: %v", errs)
	}

	mod, modok := yangE[testModule]
	if !modok {
		t.Fatalf("could not find expected module: %s (%v)", testModule, yangE)
	}

	for _, entry := range mod.Dir {
		if entry.Kind != yang.LeafEntry {
			continue
		}

		if len(entry.Errors) != 0 {
			t.Errorf("entry had associated errors: %v", entry.Errors)
			continue
		}

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
				_, gotMatch = checkPattern(ext.Argument, entry.Type.POSIXPattern)
			} else {
				// Handle unions
				results := make([]bool, 0)
				for _, membertype := range entry.Type.Type {
					// Only do the test when there is a pattern specified against the
					// type as it may not be a string.
					if membertype.Kind != yang.Ystring || len(membertype.POSIXPattern) == 0 {
						continue
					}
					matchedAllForType := true
					_, matchedAllForType = checkPattern(ext.Argument, membertype.POSIXPattern)
					results = append(results, matchedAllForType)
				}

				gotMatch = false
				for _, r := range results {
					if r == true {
						gotMatch = true
					}
				}
			}

			if gotMatch != wantMatch {
				t.Errorf("%s: string %s did not have expected result: %v", entry.Type.Name, ext.Argument, wantMatch)
			}
		}
	}
}

// checkPattern builds and compils
func checkPattern(testData string, patterns []string) (compileErr error, matched bool) {
	for _, pattern := range patterns {
		if r, err := regexp.CompilePOSIX(pattern); err != nil {
			return
		} else {
			matched = r.MatchString(testData)
		}
	}
	return
}

// init sets up the test, particularly parsing the OpenConfig path which is
// supplied as a command line argument.
func init() {
	ocdir = os.Getenv("OCDIR")
	if ocdir == "" {
		log.Fatal("missing environment variable $OCDIR for specifying model root directory")
	}
}
