package client

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ast "slice2go/ast"
	"slice2go/slicetemplate"

	log "github.com/sirupsen/logrus"
)

type ClientFlag struct {
	File       string
	Prefix     string
	Template   string
	Interface  string
	ExcludeIce bool
}

func Do(flag ClientFlag) {
	log.Debug(flag)
	ftoken, err := ast.Parse(flag.File)
	if err != nil {
		return
	}

	path := filepath.Join(flag.Prefix, strings.ToLower(ftoken.Module.Name))

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Errorf("path:%v, %v", path, err)
		return
	}

	flag.File, _ = filepath.Abs(flag.File)

	// make cpp to .a
	log.Debugf("file:%v", ftoken.Name)
	log.Debugf("module:%v", ftoken.Module.Name)
	st := slicetemplate.SliceTemplate{
		Template: flag.Template,
		Prefix:   flag.Prefix,
	}
	var m slicetemplate.Modules
	m.Module = ftoken.Module.Name
	m.ToCPP = !flag.ExcludeIce
	for n, i := range ftoken.Module.Interfaces {
		m.Interface = i.Name
		m.File = strings.TrimSuffix(ftoken.Name, ".ice")
		if len(flag.Interface) > 0 && flag.Interface != m.Interface {
			continue
		}
		log.Debugf("%v    interface:%v", n, i.Name)
		for _, f := range i.Funcs {
			log.Debugf("        function:%v", f.Name)
			if TwoBytes(f.Params) {
				m.Functions = append(m.Functions, f.Name)
			} else {
				log.Debugf("        function:%v not match two bytes", f.Name)
			}
		}
		if err := st.Make(m); err != nil {
			log.Errorf("failed to make go file: %v\n", err)
		}
	}

	if !flag.ExcludeIce {
		if err := st.MakeMakefile(m); err != nil {
			log.Errorf("failed to make MakeMakefile: %v\n", err)
		}

		cmdstr := fmt.Sprintf("cd %s && cd ice_interface && slice2cpp %s && make", path, flag.File)
		log.Debugf("cmdstr:%v", cmdstr)
		cmd := exec.Command("bash", "-c", cmdstr)

		var stdout bytes.Buffer
		cmd.Stderr = &stdout
		cmd.Stdout = &stdout
		err := cmd.Run()
		if err != nil {
			log.Errorf("failed to call tocpp: %v\n", err)
		}
		log.Debugf("make ice cpp: %s", stdout.String())
	}
}

func TwoBytes(params []interface{}) bool {
	if len(params) != 2 {
		return false
	}

	/*
		for _, v := range params {
			ps, ok := v.(*ast.TParam)
			if !ok {
				return false
			}

			if ps.PType != "bytes" {
				return false
			}
		}
	*/

	return true
}
