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

package patterncheck

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckRegexps(t *testing.T) {
	tests := []struct {
		desc             string
		inFiles          []string
		inPaths          []string
		wantFailMessages []string
	}{{
		desc:    "passing cases",
		inFiles: []string{"testdata/passing.yang"},
		inPaths: []string{"../../testdata"},
	}, {
		desc:    "simple leaf fail",
		inFiles: []string{"testdata/simple-leaf-fail.yang"},
		inPaths: []string{"../../testdata"},
		wantFailMessages: []string{
			"| `ipv-0` | `string` | `ipv4` matched but shouldn't |",
			"| `ipv-0` | `string` | `ipv6` did not match |",
		},
	}, {
		desc:    "union leaf fail",
		inFiles: []string{"testdata/union-leaf-fail.yang"},
		inPaths: []string{"../../testdata"},
		wantFailMessages: []string{
			"| `ipv-0` | `ip-string-typedef` | `ipv4` matched but shouldn't |",
			"| `ipv-0` | `ip-string-typedef` | `ipv5` did not match |",
		},
	}, {
		desc:    "derived string type fail",
		inFiles: []string{"testdata/derived-string-fail.yang"},
		inPaths: []string{"../../testdata"},
		wantFailMessages: []string{
			"| `ipv-0` | `ipv4-address-str` | `ipV4` did not match |",
			"| `ipv-0` | `ipv4-address-str` | `ipV4-address` matched but shouldn't |",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := CheckRegexps(tt.inFiles, tt.inPaths)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tt.wantFailMessages); diff != "" {
				t.Errorf("(-got, +want):\n%s", diff)
			}
		})
	}
}
