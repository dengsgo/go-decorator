package main

import (
	"encoding/json"
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var decoratorBinaryPath = os.Getenv("GOPATH") + "/bin/decorator"

type _packageInfo struct {
	Dir,
	ImportPath,
	Name,
	Target,
	Root,
	StaleReason string
	Stale  bool
	Module struct {
		Main bool
		Path,
		Dir,
		GoMod,
		GoVersion string
	}
	Match,
	GoFiles,
	Imports,
	Deps []string
}

func getPackageInfo(pkgPath string) (*_packageInfo, error) {
	command := []string{"go", "list", "-json"}
	if pkgPath != "" && pkgPath != "main" {
		command = append(command, pkgPath)
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = projectDir
	cmd.Env = os.Environ()
	bf, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	p := &_packageInfo{}
	err = json.Unmarshal(bf, p)
	if err != nil {
		return nil, err
	}
	return p, nil
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
			var name string
			pkg, _ := strconv.Unquote(ip.Path.Value)
			extName := strings.TrimRight(
				filepath.Base(pkg),
				filepath.Ext(pkg),
			)

			// e.g. path/u/name/v2
			if strings.HasPrefix(extName, "v") && func() bool {
				v, err := strconv.Atoi(strings.TrimLeft(extName, "v"))
				return err == nil && v > 1
			}() {
				arr := strings.Split(pkg, "/")
				if len(arr) > 1 {
					extName = arr[len(arr)-2]
				}
			}

			if ip.Name == nil {
				// import path/name // name form pkg
				name = extName
			} else {
				switch ip.Name.Name {
				case "":
					// import path/name // name form pkg
					name = extName
				case "_":
					// import _ path/name // name pkg, about to be replaced
					name = extName
				case ".":
					// import . path/name // ""
					name = extName
				default:
					// import yname path/name // yname from alias
					name = ip.Name.Name
				}
			}

			nameMap[name] = pkg
			pathObjMap[pkg] = ip
			pathMap[pkg] = func() string {
				if ip.Name != nil {
					return ip.Name.Name
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
