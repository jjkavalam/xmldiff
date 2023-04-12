package xmldiff

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type parser struct {
	dec     *xml.Decoder
	next    xml.Token
	nextErr error
}

func newParser(xmlData string) *parser {
	dec := xml.NewDecoder(bytes.NewBufferString(xmlData))
	return &parser{
		dec: dec,
	}
}

func (p *parser) token() (xml.Token, error) {
	t, err := p.dec.Token()

	// skip special tokens
	switch t.(type) {
	case xml.ProcInst, xml.Comment, xml.Directive:
		return p.token()
	}

	return t, err

}

func (p *parser) peek() (xml.Token, error) {
	if p.next == nil {
		p.next, p.nextErr = p.token()
	}
	return p.next, p.nextErr
}

func (p *parser) pop() (xml.Token, error) {
	t, err := p.next, p.nextErr
	if t != nil {
		p.next = nil
		p.nextErr = nil
		return t, err
	}
	return p.token()
}

// parseTag converts a stream of tokens into a tree structure,
// using a parsing technique called "recursive descent parsing".
// Each parsing function takes a stream of tokens as input and return a tree structure.
func (p *parser) parseTag() (*Tag, error) {

	tg := &Tag{
		Children: make([]*Tag, 0),
	}

	err := p.takeWhiteSpace()

	if err != nil {
		return nil, err
	}

	t, err := p.pop()

	if err != nil {
		return nil, err
	}

	if tt, ok := t.(xml.StartElement); !ok {
		var s string
		if tt, ok := t.(xml.CharData); ok {
			s = string(tt)
		}
		return nil, fmt.Errorf("expected start element; found %T '%s'", t, s)
	} else {
		tg.Name = tt.Name.Local
	}

	err = p.takeWhiteSpace()

	if err != nil {
		return nil, err
	}

	t, err = p.peek()

	if err != nil {
		return nil, err
	}

	switch tt := t.(type) {
	case xml.CharData:
		_, _ = p.pop()
		tg.Value = string(tt)
	case xml.EndElement:
	default:
		for {
			ct, err := p.peek()
			if err != nil {
				return nil, err
			}
			if _, ok := ct.(xml.EndElement); ok {
				// exit loop parsing child elements on finding an unmatched end element
				break
			}
			ctg, err := p.parseTag()
			if err != nil {
				return nil, err
			}
			tg.Children = append(tg.Children, ctg)
		}
	}

	t, err = p.pop()

	if err != nil {
		return nil, err
	}

	if _, ok := t.(xml.EndElement); !ok {
		return nil, fmt.Errorf("expected end element; found %#v", t)
	}

	err = p.takeWhiteSpace()
	if err != nil && err != io.EOF {
		return nil, err
	}

	return tg, nil
}

// takeWhitespace checks if next token is an empty whitespace
// it consumes that token and returns if so; otherwise does not consume the token.
func (p *parser) takeWhiteSpace() error {
	// keep consuming tokens as long as they are empty whitespaces
	// why would two whitespace tokens come one after the another rather than as one combined token ?
	// Because, there could be other "special" tokens in between ! (see p.token() for details)
	for {
		t, err := p.peek()

		if err != nil {
			return err
		}

		if t, ok := t.(xml.CharData); ok {
			s := string(t)
			if strings.TrimSpace(s) == "" {
				_, _ = p.pop()
			} else {
				break
			}
		} else {
			break
		}
	}
	return nil
}
