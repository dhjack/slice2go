package slicetemplate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func ToGoFunName(a string) string {
	if len(a) > 0 {
		return strings.ToUpper(a[0:1]) + a[1:]
	}
	return a
}

var (
	_funcMap = template.FuncMap{
		"tolower":     strings.ToLower,
		"togofunname": ToGoFunName,
	}
)

type SliceTemplate struct {
	Template   string
	Prefix     string
	filePrefix string
}

type Modules struct {
	File      string
	Module    string
	Interface string
	Functions []string
	ToCPP     bool
}

func (st *SliceTemplate) Make(m Modules) error {
	os.MkdirAll(filepath.Join(st.Prefix, strings.ToLower(fmt.Sprintf("%s/%s", m.Module, m.Interface))), os.ModePerm)

	st.filePrefix = filepath.Join(st.Prefix, strings.ToLower(fmt.Sprintf("%s/%s/%s", m.Module, m.Interface, m.Interface)))
	if err := st.makeh(m); err != nil {
		return err
	}
	if err := st.makeCC(m); err != nil {
		return err
	}
	if err := st.makeGo(m); err != nil {
		return err
	}
	return nil
}

func (st *SliceTemplate) MakeMakefile(m Modules) error {
	t, err := st.makeTemplate("module_makefile.tmpl", module_makefile)
	if err != nil {
		log.Println("ParseFiles err:", err)
		return err
	}

	os.MkdirAll(filepath.Join(st.Prefix, strings.ToLower(fmt.Sprintf("%s/ice_interface", m.Module))), os.ModePerm)
	outputFile, outputError := os.OpenFile(filepath.Join(st.Prefix, fmt.Sprintf("%s/ice_interface/Makefile", strings.ToLower(m.Module))), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if outputError != nil {
		log.Println(outputError)
		return outputError
	}
	defer outputFile.Close()

	err = t.Execute(outputFile, m)
	if err != nil {
		log.Println("executing template:", err)
		return err
	}

	return nil
}

func (st *SliceTemplate) makeh(m Modules) error {
	t, err := st.makeTemplate("module_wrap_h.tmpl", module_h)
	if err != nil {
		log.Println("ParseFiles err:", err)
		return err
	}

	outputFile, outputError := os.OpenFile(st.filePrefix+"_wrap.h", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if outputError != nil {
		log.Println(outputError)
		return outputError
	}
	defer outputFile.Close()

	err = t.Execute(outputFile, m)
	if err != nil {
		log.Println("executing template:", err)
		return err
	}

	return nil
}

func (st *SliceTemplate) makeTemplate(name, str string) (*template.Template, error) {
	if len(st.Template) > 0 {
		return template.New(name).Funcs(_funcMap).ParseFiles(filepath.Join(st.Template, name))
	}
	return template.Must(template.New(name).Funcs(_funcMap).Parse(source(str))), nil
}

func (st *SliceTemplate) makeCC(m Modules) error {
	t, err := st.makeTemplate("module_wrap_cc.tmpl", module_cpp)
	if err != nil {
		log.Println("ParseFiles err:", err)
		return err
	}
	outputFile, outputError := os.OpenFile(st.filePrefix+"_wrap.cc", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if outputError != nil {
		log.Println(outputError)
		return outputError
	}
	defer outputFile.Close()

	err = t.Execute(outputFile, m)
	if err != nil {
		log.Println("executing template:", err)
		return err
	}

	return nil
}

func (st *SliceTemplate) makeGo(m Modules) error {
	t, err := st.makeTemplate("module.go.tmpl", module_go)
	if err != nil {
		log.Println("ParseFiles err:", err)
		return err
	}
	outputFile, outputError := os.OpenFile(st.filePrefix+".go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if outputError != nil {
		log.Println(outputError)
		return outputError
	}
	defer outputFile.Close()

	err = t.Execute(outputFile, m)
	if err != nil {
		log.Println("executing template:", err)
		return err
	}

	return nil
}
