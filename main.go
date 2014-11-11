package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var header string

var DefaultExcludes = []string{
	// miscellaneous typical temporary files
	"**/*~",
	"**/#*#",
	"**/.#*",
	"**/%*%",
	"**/._*",
	"**/.repository/**",

	// CVS
	"**/CVS",
	"**/CVS/**",
	"**/.cvsignore",

	// RCS
	"**/RCS",
	"**/RCS/**",

	// SCCS
	"**/SCCS",
	"**/SCCS/**",

	// Visual SourceSafe
	"**/vssver.scc",

	// Subversion
	"**/.svn",
	"**/.svn/**",

	// Arch
	"**/.arch-ids",
	"**/.arch-ids/**",

	// Bazaar
	"**/.bzr",
	"**/.bzr/**",

	//SurroundSCM
	"**/.MySCMServerInfo",

	// Mac
	"**/.DS_Store",

	// Serena Dimensions Version 10
	"**/.metadata",
	"**/.metadata/**",

	// Mercurial
	"**/.hg",
	"**/.hg/**",
	"**/.hgignore",

	// git
	"**/.git",
	"**/.git/**",
	"**/.gitignore",
	"**/.gitmodules",

	// BitKeeper
	"**/BitKeeper",
	"**/BitKeeper/**",
	"**/ChangeSet",
	"**/ChangeSet/**",

	// darcs
	"**/_darcs",
	"**/_darcs/**",
	"**/.darcsrepo",
	"**/.darcsrepo/**",
	"**/-darcs-backup*",
	"**/.darcs-temp-mail",

	// maven project's temporary files
	"**/target/**",
	"**/test-output/**",
	"**/release.properties",
	"**/dependency-reduced-pom.xml",
	"**/pom.xml.releaseBackup",

	// code coverage tools
	"**/cobertura.ser",
	"**/.clover/**",

	// eclipse project files
	"**/.classpath",
	"**/.project",
	"**/.settings/**",

	// IDEA projet files
	"**/*.iml",
	"**/*.ipr",
	"**/*.iws",
	".idea/**",

	// descriptors
	"**/MANIFEST.MF",

	// binary files - images
	"**/*.jpg",
	"**/*.png",
	"**/*.gif",
	"**/*.ico",
	"**/*.bmp",
	"**/*.tiff",
	"**/*.tif",
	"**/*.cr2",
	"**/*.xcf",

	// binary files - programs
	"**/*.class",
	"**/*.exe",
	"**/*.dll",
	"**/*.so",

	// checksum files
	"**/*.md5",
	"**/*.sha1",

	// binary files - archives
	"**/*.jar",
	"**/*.zip",
	"**/*.rar",
	"**/*.tar",
	"**/*.tar.gz",
	"**/*.tar.bz2",
	"**/*.gz",

	// binary files - documents
	"**/*.xls",

	// ServiceLoader files
	"**/META-INF/services/**",

	// Markdown files
	"**/*.md",

	// Office documents
	"**/*.xls",
	"**/*.doc",
	"**/*.odt",
	"**/*.ods",
	"**/*.pdf",

	// Travis
	"**/.travis.yml",

	// flash
	"**/*.swf",

	// json files
	"**/*.json",
}

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

			log.Println(fpath)

			buf, err := ioutil.ReadFile(fpath)
			checkErr("Reads file error: ", err)

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
	header = buf.String()

	root := &FileNode{Path: dirPath, Type: "d", FileNodes: []*FileNode{}}
	walk(dirPath, root)

}

func checkErr(errMsg string, err error) {
	if nil != err {
		log.Fatal(errMsg, err)
	}
}
