package ast

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type TFile struct {
	KindType string
	Name     string
	Module   *TModule
}

type TModule struct {
	KindType   string
	Name       string
	Interfaces []*TInterface
}

type TInterface struct {
	KindType string
	Name     string
	Funcs    []*TFunction
}

type TFunction struct {
	KindType    string
	Name        string
	Params      []interface{}
	ReturnValue interface{}
}

type TParam struct {
	KindType string
	Name     string
	PType    string
}

func Parse(path string) (*TFile, error) {
	t, err := NewTokener(path)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var file TFile
	file.KindType = "file"
	file.Name = filepath.Base(path)
	file.parse(t)
	return &file, nil
}

func (f *TFile) parse(t *Tokener) {
	for {
		token, ok := t.Next()
		if !ok {
			return
		}
		//fmt.Println(token)
		switch token {
		case "module":
			m := new(TModule)
			m.KindType = "module"
			f.Module = m
			m.parse(t)
		case "/":
			SkipComment(t)
		default:
		}
	}
}

func (m *TModule) parse(t *Tokener) {
	leftCurlyBrackets, rightCurlyBrackets := false, false
	leftSquareBrackets, rightSquareBrackets := false, false

	for {
		token, ok := t.Next()
		if !ok {
			return
		}
		//fmt.Println(token)
		switch token {
		case "interface":
			ti := new(TInterface)
			ti.KindType = "interface"
			m.Interfaces = append(m.Interfaces, ti)
			ti.parse(t)
		case "\n":
		case "{":
			leftCurlyBrackets = true
		case "}":
			rightCurlyBrackets = true
		case "/":
			SkipComment(t)
		default:
			if !leftCurlyBrackets {
				m.Name = token
			}

			if leftCurlyBrackets && rightCurlyBrackets {
			}

			if leftSquareBrackets && rightSquareBrackets {
			}
		}
	}
}

func (ii *TInterface) parse(t *Tokener) {
	leftCurlyBrackets, rightCurlyBrackets := false, false
	leftSquareBrackets, rightSquareBrackets := false, false

	for {
		token, ok := t.Next()
		if !ok {
			return
		}
		//fmt.Println(token)
		switch token {
		case "{":
			leftCurlyBrackets = true
		case "}":
			rightCurlyBrackets = true
		case "[":
			leftSquareBrackets = true
		case "]":
			rightSquareBrackets = true
		case ";":
			return
		case "\n":
		case "/":
			SkipComment(t)
		case "idempotent":
		default:
			if !leftCurlyBrackets {
				if len(ii.Name) == 0 {
					ii.Name = token
				}
				continue
			}
			if leftSquareBrackets && !rightSquareBrackets {
				continue
			}

			if leftCurlyBrackets && !rightCurlyBrackets {
				tf := new(TFunction)
				tf.KindType = "funcation"
				ii.Funcs = append(ii.Funcs, tf)
				tf.ReturnValue = token
				tf.parse(t)
			}
		}
	}
}

func (tf *TFunction) parse(t *Tokener) {
	leftCurlyBrackets, rightCurlyBrackets := false, false
	leftSquareBrackets, rightSquareBrackets := false, false
	leftParentheses, rightParentheses := false, false

	for {
		token, ok := t.Next()
		if !ok {
			return
		}
		//fmt.Println(token)
		switch token {
		case "{":
			leftCurlyBrackets = true
		case "}":
			rightCurlyBrackets = true
		case "[":
			leftSquareBrackets = true
		case "]":
			rightSquareBrackets = true
		case "(":
			leftParentheses = true
		case ")":
			rightParentheses = true
		case "/":
			SkipComment(t)
		case ";":
			return
		default:
			if !leftParentheses {
				tf.Name = token
				continue
			}
			if leftParentheses && !rightParentheses {
				tp := new(TParam)
				tp.KindType = "param"
				tf.Params = append(tf.Params, tp)
				t.PushToken(token)
				tp.parse(t)
			}

			if leftCurlyBrackets && !rightCurlyBrackets {
			}

			if leftSquareBrackets && rightSquareBrackets {
			}
		}
	}
}

func (tp *TParam) parse(t *Tokener) {
	// Square []
	// Curly {}
	// Parentheses ()
	for {
		token, ok := t.Next()
		if !ok {
			return
		}
		//fmt.Println(token)
		switch token {
		case ")":
			t.PushToken(token)
			return
		case "/":
			SkipComment(t)
		case ",":
			return
		case "bytes":
			tp.PType = "bytes"
		case "out":
		default:
			if len(tp.Name) > 0 {
				tp.PType = tp.Name
			}
			tp.Name = token
		}
	}
	return
}

func SkipComment(t *Tokener) bool {
	preToken := "/"
	token, ok := t.Next()
	if !ok {
		return false
	}

	//fmt.Println("/comment1:", token)
	isCommentOne := false
	isCommentSecond := false

	if token == "/" {
		isCommentOne = true
	} else if token == "*" {
		isCommentSecond = true
	} else {
		t.PushToken(token)
		return false
	}

	for {
		token, ok := t.Next()
		if !ok {
			return false
		}
		//fmt.Println("comment2:", token)
		switch token {
		case "/":
			if isCommentSecond && preToken == "*" {
				return true
			}
		case "*":
		case "\n":
			if isCommentOne {
				return true
			}
		}
		preToken = token
	}
	return true
}
