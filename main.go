package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
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

			buf, err := ioutil.ReadFile(fpath)
			checkErr("Reads file error: ", err)

			ext := filepath.Ext(fpath)

			handler := GetHandler(ext)
			if nil == handler {
				log.Printf("Can't handle [%s]", ext)

				continue
			}

			originalContent := string(buf)
			header := handler.Execute(rawHeader)

			action := getAction(originalContent, header)
			var content string

			switch action {
			case "no":
				continue
			case "add":
				content = header + "\n" + originalContent
				defer log.Printf("Added header to file [%s]", fpath)
			case "update":
				headerLines := strings.Split(header, "\n")
				originalContentLines := strings.Split(originalContent, "\n")
				contentLines := originalContentLines[len(headerLines):]
				headerLines = append(headerLines, contentLines...)

				content = strings.Join(headerLines, "\n")
				defer log.Printf("Updated header to file [%s]", fpath)
			default:
				log.Fatalf("Wrong action [%s]", action)
			}

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

// getAction returns a handle action for the specified original content and header.
//
//  1. "add" means need to add the header to the original content
//  2. "update" means need to update (replace) the header of the original content
//  3. "no" means nothing need to do
func getAction(originalContent, header string) string {
	headerLines := strings.Split(header, "\n")
	originalContentLines := strings.Split(originalContent, "\n")

	if len(headerLines) > len(originalContentLines) {
		return "add"
	}

	originalHeaderLines := originalContentLines[:len(headerLines)]

	result := similar(originalHeaderLines, headerLines)
	if 100 <= result {
		return "no"
	}

	if result >= 70 {
		return "update"
	}

	return "no"
}

// [0, 100]
//
//  0: not similar at all
//  100: as the same
func similar(lines1, lines2 []string) int {
	if len(lines1) != len(lines2) {
		return 0
	}

	length := len(lines1)
	same := 0
	for i := 0; i < length; i++ {
		l1 := strings.TrimSpace(lines1[i])
		l2 := strings.TrimSpace(lines2[i])

		if l1 == l2 {
			same++
		}
	}

	return int(math.Floor(float64(same) / float64(length) * 100))
}

func main() {
	headerTemplate := "test/test_header.txt"

	dir := "."

	includes := []string{"test/test_*.go"}
	excludes := []string{""}

	useDefaultExcludes := true

	mapping := map[string]string{"java": "SLASHSTAR_STYLE"}
	// TODO: default mapping

	properties := map[string]string{"year": "2014", "owner": "Liang Ding"}

	// TODO: mapping
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
		log.Fatal(errMsg+", caused by: ", err)
	}
}
