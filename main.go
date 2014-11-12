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

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Conf struct {
	Dir                string
	Template           string
	Includes           []string
	Excludes           []string
	UseDefaultExcludes bool
	Properties         map[string]string
}

var conf = &Conf{}

var rawHeader *RawHeader

type FileNode struct {
	Path      string
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
			walk(fpath, &child)
		} else {
			match := match(fpath)

			if !match {
				continue
			}

			buf, err := ioutil.ReadFile(fpath)
			checkErr("Reads file error", err)

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

	for _, exclude := range conf.Excludes {
		if AntPathMatch(exclude, path) {
			return false
		}
	}

	for _, include := range conf.Includes {
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

	return "add"
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
	bts, err := ioutil.ReadFile(".header.json")
	checkErr("Loads configuration error", err)

	err = json.Unmarshal(bts, conf)
	checkErr("Parses configuration error", err)

	dir := conf.Dir
	properties := conf.Properties

	dirPath, err := filepath.Abs(dir)
	checkErr("Reads dir path error", err)

	// Includes
	for _, include := range conf.Includes {
		path := filepath.Join(dirPath, include)
		path = filepath.ToSlash(path)
		conf.Includes = append(conf.Includes, path)
	}

	// Excludes
	if conf.UseDefaultExcludes {
		conf.Excludes = append(conf.Excludes, DefaultExcludes...)
	}
	for _, exclude := range conf.Excludes {
		path := filepath.Join(dirPath, exclude)
		path = filepath.ToSlash(path)
		conf.Excludes = append(conf.Excludes, path)
	}

	t, err := template.ParseFiles(conf.Template)
	checkErr("Can't find header template ["+conf.Template+"]", err)

	var buf bytes.Buffer
	t.Execute(&buf, properties)
	rawHeader = NewRawHeader(buf.String())

	root := &FileNode{Path: dirPath, FileNodes: []*FileNode{}}
	walk(dirPath, root)

}

func checkErr(errMsg string, err error) {
	if nil != err {
		log.Fatal(errMsg+", caused by: ", err)
	}
}
