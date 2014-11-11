package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestMatch(t *testing.T) {
	match := AntPathMatch("/home/**/*.go", "/home/daniel/test/a.go")
	fmt.Println(match)

	pattern := filepath.ToSlash("**\\header\\*.go")
	path := filepath.ToSlash("D:\\GoGoGo\\src\\github.com\\88250\\header\\main.go")
	match = AntPathMatch(pattern, path)
	fmt.Println(match)
}
