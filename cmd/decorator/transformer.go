package main

import (
	"bytes"
	"fmt"
	"github.com/dengsgo/go-decorator/cmd/logs"
	"os"
	"os/exec"
	"path"
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
	cmd.Env = cmd.Environ()
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
	q   map[string]string
	cfg string

	qMap map[string]bool
}

func newImporter(cfg string) *importer {
	return &importer{
		q:    map[string]string{},
		cfg:  cfg,
		qMap: map[string]bool{},
	}
}

func (i *importer) addImport(pkg string) (err error) {
	if _, ok := i.qMap[pkg]; ok {
		return
	}
	i.qMap[pkg] = true
	pg := getPkgCompiledInfo(pkg)
	i.q[pkg] = fmt.Sprintf("packagefile %s=%s\n", pkg, pg.export)
	return
}

func (i *importer) sync() (err error) {
	if len(i.q) == 0 {
		return
	}

	bs, err := os.ReadFile(i.cfg)
	if err != nil {
		return
	}
	bf := bytes.NewBuffer(bs)
	for {
		line, err := bf.ReadString('\n')
		if err != nil {
			break
		}
		if !strings.HasPrefix(line, "packagefile ") {
			continue
		}
		kv := line[len("packagefile "):]
		index := strings.Index(kv, "=")
		if index < 0 {
			continue
		}
		p := strings.TrimSpace(kv[0:index])
		if _, ok := i.q[p]; ok {
			delete(i.q, p)
		}
	}
	if len(i.q) == 0 {
		return
	}
	insert := bytes.NewBuffer([]byte{})
	for _, v := range i.q {
		insert.WriteString(v)
	}
	bs = append(bs, insert.Bytes()...)
	logs.Debug("added file:", string(bs))
	err = os.WriteFile(i.cfg, bs, 0666)
	appendPackagefileSharedFile(insert.Bytes())
	i.q = make(map[string]string)
	return
}

func appendPackagefileSharedFile(bs []byte) {
	fileName := path.Join(tempDir, "sharedPackagefile.txt")
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logs.Error("appendPackagefileSharedFile OpenFile fail", fileName, err, string(bs))
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err := file.Write(bs); err != nil {
		logs.Error("appendPackagefileSharedFile Write fail", fileName, err, string(bs))
	}
	if file.Sync() != nil {
		logs.Error("appendPackagefileSharedFile Sync fail", fileName, err, string(bs))
	}
}
