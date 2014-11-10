package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	headerTemplate := "test-header.txt"

	includes := []string{"**/xxx/*.go"}

	excludes := []string{"**/xxx/a.go"}

	useDefaultExcludes := true
	// TODO: default excludes

	mapping := map[string]string{"java": "SLASHSTAR_STYLE"}
	// TODO: default mapping

	properties := map[string]string{"year": "2014", "owner": "Liang Ding"}

	// TODO: includes, excludes, useDefaultExcludes, mapping
	_ = includes
	_ = excludes
	_ = useDefaultExcludes
	_ = mapping

	t, err := template.ParseFiles(headerTemplate)

	if nil != err {
		log.Fatalf("Can't find header template [%s]", headerTemplate)

		return
	}

	t.Execute(os.Stdout, properties)
}
