package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

var header string

type FileNode struct {
	Path      string
	Type      string // "f": file, "d": directory
	FileNodes []*FileNode
}

func walk(path string, node *FileNode) {
	files := listFiles(path)

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

			if !match(fpath, "TODO: pattern") {
				continue
			}

			buf, err := ioutil.ReadFile(fpath)

			checkErr("Reads file error: ", err)

			log.Println(fpath)

			//content := header + string(buf)

			//fout, err := os.Create(fpath)
			//checkErr("Prepare to write file err: ", err)

			//fout.WriteString(content)

			//err = fout.Close()
			//checkErr("Write file err: ", err)

			_ = buf
		}
	}
}

func listFiles(dirname string) []string {
	f, _ := os.Open(dirname)

	names, _ := f.Readdirnames(-1)
	f.Close()

	sort.Strings(names)

	ret := []string{}

	for _, name := range names {
		fio, _ := os.Lstat(filepath.Join(dirname, name))

		// TODO: excludes
		if ".git" == fio.Name() {
			continue
		}

		ret = append(ret, name)
	}

	return ret
}

func match(path, pattern string) bool {
	return strings.HasSuffix(path, ".go")
}

func main() {
	headerTemplate := "test-header.txt"

	dir := "."

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

	dirPath, err := filepath.Abs(dir)
	checkErr("Reads dir path error: ", err)

	t, err := template.ParseFiles(headerTemplate)
	checkErr("Can't find header template ["+headerTemplate+"]", err)

	var buf bytes.Buffer
	t.Execute(&buf, properties)
	header = buf.String()

	root := &FileNode{Path: dirPath, Type: "d", FileNodes: []*FileNode{}}
	walk(dirPath, root)

}

func checkErr(errMsg string, err error) {
	if nil != err {
		log.Fatal(errMsg, err)
	}
}
