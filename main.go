package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var rawHeader *RawHeader

var Includes = []string{}
var Excludes = []string{}

type FileNode struct {
	Path      string
	Type      string // "f": file, "d": directory
	FileNodes []*FileNode
}

func walk(path string, node *FileNode) {
	dirFile, _ := os.Open(path)

	files, _ := dirFile.Readdirnames(-1)
	dirFile.Close()

	for _, filename := range files {
		fpath := filepath.Join(path, filename)

		fio, err := os.Lstat(fpath)

		checkErr("", err)

		child := FileNode{Path: fpath, FileNodes: []*FileNode{}}
		node.FileNodes = append(node.FileNodes, &child)

		if fio.IsDir() {
			child.Type = "d"

			walk(fpath, &child)
		} else {
			child.Type = "f"

			match := match(fpath)

			if !match {
				continue
			}

			log.Println("path", fpath)

			buf, err := ioutil.ReadFile(fpath)
			checkErr("Reads file error: ", err)

			ext := filepath.Ext(fpath)

			handler := GetHandler(ext)
			if nil == handler {
				log.Printf("Can't handle [%s]", ext)

				continue
			}

			header := handler.Execute(rawHeader)
			content := header + string(buf)

			fout, err := os.Create(fpath)
			checkErr("Prepare to write file err: ", err)

			fout.WriteString(content)

			err = fout.Close()
			checkErr("Write file err: ", err)

		}
	}
}

func match(path string) bool {
	path = filepath.ToSlash(path)

	for _, exclude := range Excludes {
		if AntPathMatch(exclude, path) {
			return false
		}
	}

	for _, include := range Includes {
		if AntPathMatch(include, path) {
			return true
		}

	}

	return false
}

func main() {
	headerTemplate := "test-header.txt"

	dir := "."

	includes := []string{"*.go"}
	excludes := []string{""}

	useDefaultExcludes := true

	mapping := map[string]string{"java": "SLASHSTAR_STYLE"}
	// TODO: default mapping

	properties := map[string]string{"year": "2014", "owner": "Liang Ding"}

	// TODO: includes, excludes, useDefaultExcludes, mapping
	_ = includes
	_ = excludes
	_ = useDefaultExcludes
	_ = mapping

	dirPath, err := filepath.Abs(dir)
	checkErr("Reads dir path error: ", err)

	// Includes
	for _, include := range includes {
		path := filepath.Join(dirPath, include)
		path = filepath.ToSlash(path)
		Includes = append(Includes, path)
	}

	// Excludes
	if useDefaultExcludes {
		excludes = append(excludes, DefaultExcludes...)
	}
	for _, exclude := range excludes {
		path := filepath.Join(dirPath, exclude)
		path = filepath.ToSlash(path)
		Excludes = append(Excludes, path)
	}

	t, err := template.ParseFiles(headerTemplate)
	checkErr("Can't find header template ["+headerTemplate+"]", err)

	var buf bytes.Buffer
	t.Execute(&buf, properties)
	rawHeader = NewRawHeader(buf.String())

	root := &FileNode{Path: dirPath, Type: "d", FileNodes: []*FileNode{}}
	walk(dirPath, root)

}

func checkErr(errMsg string, err error) {
	if nil != err {
		log.Fatal(errMsg, err)
	}
}
