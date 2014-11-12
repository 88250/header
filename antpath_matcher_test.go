// Copyright (c) 2014, B3log
//  
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//  
//     http://www.apache.org/licenses/LICENSE-2.0
//  
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import "testing"

type testcase struct {
	pattern string
	path    string
	match   bool
}

var testcases = []testcase{
	testcase{pattern: "/home/daniel/test/a.go", path: "/home/daniel/test/a.go", match: true},
	testcase{pattern: "/home/daniel/test/a.go", path: "/home/daniel/test/a1.go", match: false},

	testcase{pattern: "/home/daniel/test/*.go", path: "/home/daniel/test/a.go", match: true},
	testcase{pattern: "/home/daniel/test/*.go", path: "/home/daniel/test/ab.go", match: true},
	testcase{pattern: "/home/daniel/test/?.go", path: "/home/daniel/test/a.go", match: true},
	testcase{pattern: "/home/daniel/test/?.go", path: "/home/daniel/test/ab.go", match: false},
	testcase{pattern: "/home/daniel/test/*.go", path: "/home/daniel/test/a.go1", match: false},

	testcase{pattern: "/home/**/*.go", path: "/home/daniel/test/a.go", match: true},
	testcase{pattern: "/home/**/*.go", path: "/home/daniel/test/a.go1", match: false},

	testcase{pattern: "/home/?/a.go", path: "/home/d/a.go", match: true},
	testcase{pattern: "/home/?/*.go", path: "/home/d/a.go", match: true},
	testcase{pattern: "/home/?/*.go", path: "/home/d/a.go1", match: false},
	testcase{pattern: "/home/?/a.go", path: "/home/d1/a.go", match: false},
}

func TestMatch(t *testing.T) {
	for _, tc := range testcases {
		match := AntPathMatch(tc.pattern, tc.path)

		if match != tc.match {
			t.Fatalf("pattern [%s], path [%s], result [%v], expected [%v]", tc.pattern, tc.path, match, tc.match)
		}
	}
}
