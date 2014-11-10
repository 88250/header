package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	headerTemplate := "test-header.txt"

	t, err := template.ParseFiles(headerTemplate)

	if nil != err {
		log.Fatalf("Can't find header template [%s]", headerTemplate)

		return
	}

	model := map[string]string{"year": "2014", "owner": "Liang Ding"}
	t.Execute(os.Stdout, model)
}
