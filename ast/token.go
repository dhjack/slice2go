package ast

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"
)

type Tokener struct {
	reader  io.Reader
	scanner *bufio.Scanner
	preToken string
}

func NewTokener(file string) (*Tokener, error) {
	sf, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(sf)
	scanner.Split(ScanWords)
	return &Tokener{
		reader:  sf,
		scanner: scanner,
	}, nil
}

func (t *Tokener) Next() (string, bool) {
	if len(t.preToken) > 0 {
		token := t.preToken
		t.preToken = ""
		return token, true
	}
	if t.scanner.Scan() {
		return t.scanner.Text(), true
	}
	return "", false
}

func (t *Tokener) PushToken(token string) {
	t.preToken = token
}

func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

func isSplit(r rune) bool {
	switch r {
	case '{', '}', '[', ']', ';', ',', '(', ')', '"', '\n', '*', '/':
		return true
	}
	return false
}

func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	//fmt.Printf("b>%s<b, %v\n", string(data), atEOF)
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			if isSplit(r) {
				return start + width, data[start : start+width], nil
			}
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSplit(r) {
			return i, data[start:i], nil
		}
		if isSpace(r) {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}
