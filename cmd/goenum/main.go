package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	help := flag.Bool("help", true, "Help")
	path := flag.String("path", "", "path")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	if *help && *path == "" {
		fmt.Println("goenum --path <PATH> (--debug)")
		return
	}

	if *path == "" {
		log.Fatal("goenum: Path is a required parameter; --path <PATH>")
	}

	generate(*path, *debug)
}

func generate(p string, debug bool) {
	funcs := template.FuncMap{"join": strings.Join}
	t, _ := template.New("").Funcs(funcs).Parse(tmpl)

	fileSet := token.NewFileSet()
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".go") {
			f, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
			if err == nil {
				for _, c := range f.Comments {
					for _, c1 := range c.List {
						if c1.Text == "// x-enum" {
							dir, _ := filepath.Split(path)
							d, err := parser.ParseDir(fileSet, dir, nil, parser.PackageClauseOnly)
							if err != nil {
								return err
							}
							pkgName := ""
							for k, _ := range d {
								pkgName = k
							}

							typeName := ""
							m := make(map[int]string)
							for k, v := range f.Scope.Objects {
								if v.Kind == ast.Typ {
									typeName = k
								}
								if v.Kind == ast.Con {
									m[v.Data.(int)] = k
								}
							}

							tmplModel := tmplModel{
								PkgName:  pkgName,
								Enums:    make([]EnumModel, len(m)),
								TypeName: typeName,
							}

							for k, v := range m {
								tmplModel.Enums[k] = EnumModel{
									Word:     v,
									TypeName: typeName,
								}
							}

							newFile, _ := os.Create(dir + strings.ToLower(typeName) + ".gen.go")
							w := bufio.NewWriter(newFile)
							err = t.Execute(w, tmplModel)
							if err != nil {
								fmt.Println(err)
							}
							w.Flush()

							if debug {
								fmt.Println("goenum: " + path + ": wrote " + dir + pkgName + ".gen.go")
							}
						}
					}
				}
			}
		}

		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("goenum: wrote all enum helpers")
}

type EnumModel struct {
	Word     string
	TypeName string
}

type tmplModel struct {
	PkgName  string
	Enums    []EnumModel
	TypeName string
}

const tmpl = `// Generated code. DO NOT EDIT.
package {{ .PkgName }}

import "errors"

func (kind {{ .TypeName }}) String() string {
	switch kind { {{ range $index, $element := .Enums }}
	case {{ $index }}:
		return "{{$element.Word}}" {{ end }}
	default:
		return ""
	}
}

func Parse(name string) ({{ .TypeName }}, error) {
	switch name { {{ range $index, $element := .Enums }}
	case "{{$element.Word}}":
		return {{ $element.TypeName }}({{ $index }}), nil {{ end }}
	default:
		return {{ .TypeName }}(0), errors.New("Enum for \"{{ .TypeName }}\" not found using name = " + name)
	}
}`
