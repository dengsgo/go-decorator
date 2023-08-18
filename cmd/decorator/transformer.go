package main

import (
	"bytes"
	"github.com/dengsgo/go-decorator/cmd/logs"
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const listFormat = `GO.LIST.EXPORT={{.Export}}
GO.LIST.DIR={{.Dir}}`

var decoratorBinaryPath = os.Getenv("GOPATH") + "/bin/decorator"

type pkgCompiled struct {
	work,
	export,
	dir string
}

func getPkgCompiledInfo(pkg string) *pkgCompiled {
	return pkgInfo(runGoListCommend(pkg))
}

func runGoListCommend(pkg string) *bytes.Buffer {
	logs.Debug(decoratorBinaryPath)
	var buf = bytes.NewBuffer([]byte{})
	cmd := exec.Command("go", "list", "-f", listFormat, "-export", "-toolexec", decoratorBinaryPath /*"-work",*/, pkg)
	logs.Debug("runGoListCommend", strings.Join(cmd.Args, " "))
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Dir = projectDir
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		logs.Error("runGoListCommend fail", cmd.Args, err)
	}
	logs.Debug(projectDir+"/runGoListCommend.log", buf.String())
	return buf
}

func getGoModPath() string {
	var buf = bytes.NewBuffer([]byte{})
	cmd := exec.Command("go", "list", "-f", "{{.Module.Path}}")
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Dir = projectDir
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		logs.Error("getGoModPath()", err)
	}
	return strings.TrimSpace(buf.String())
}

func pkgInfo(buf *bytes.Buffer) *pkgCompiled {
	pc := &pkgCompiled{}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(line, "WORK=") {
			pc.work = line[len("WORK=") : len(line)-1]
		} else if strings.HasPrefix(line, "GO.LIST.EXPORT=") {
			pc.export = line[len("GO.LIST.EXPORT=") : len(line)-1]
		} else if strings.HasPrefix(line, "GO.LIST.DIR=") {
			pc.dir = line[len("GO.LIST.DIR=") : len(line)-1]
		}
	}
	logs.Debug("pkgInfo", pc)
	return pc
}

type importer struct {
	nameMap    map[string]string
	pathMap    map[string]string
	pathObjMap map[string]*ast.ImportSpec
}

func newImporter(f *ast.File) *importer {
	nameMap := map[string]string{}
	pathMap := map[string]string{}
	pathObjMap := map[string]*ast.ImportSpec{}
	if f.Imports != nil && len(f.Imports) > 0 {
		for _, ip := range f.Imports {
			if ip == nil {
				continue
			}
			name := ""
			pkg, _ := strconv.Unquote(ip.Path.Value)
			if ip.Name != nil && ip.Name.Name != "" && ip.Name.Name != "_" {
				name = ip.Name.Name
			} else {
				name = filepath.Base(pkg)
			}
			nameMap[name] = pkg
			pathObjMap[pkg] = ip
			pathMap[pkg] = func() string {
				if ip.Name != nil && ip.Name.Name == "_" {
					return "_"
				}
				return name
			}()
		}
	}
	return &importer{
		nameMap:    nameMap,
		pathMap:    pathMap,
		pathObjMap: pathObjMap,
	}
}

func (i *importer) importedName(name string) (pat string, ok bool) {
	pat, ok = i.nameMap[name]
	return
}

func (i *importer) importedPath(pkg string) (name string, ok bool) {
	name, ok = i.pathMap[pkg]
	return
}
