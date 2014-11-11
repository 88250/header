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
