package xmldiff

import (
	"fmt"
	"io"
)

func Parse(xmlData string) (*Tag, error) {
	p := NewParser(xmlData)
	return p.ParseTag()
}

type Tag struct {
	Name     string
	Children []*Tag
	// Value is valid only if Children is empty
	Value string
}

func (tg *Tag) String(w io.StringWriter) error {
	_, err := w.WriteString(fmt.Sprintf("<%s>", tg.Name))
	if err != nil {
		return err
	}
	if len(tg.Children) == 0 {
		_, err = w.WriteString(tg.Value)
		if err != nil {
			return err
		}
	} else {
		for _, ctg := range tg.Children {
			err = ctg.String(w)
			if err != nil {
				return err
			}
		}
	}
	_, err = w.WriteString(fmt.Sprintf("</%s>", tg.Name))
	if err != nil {
		return err
	}
	return nil
}

// Diff compares this tag with another.
// It performs a tree traversal on both trees simultaneously and returns a list of differences between the trees.
func (tg *Tag) Diff(other *Tag, w io.StringWriter) (bool, error) {
	s := NewStack()
	s.Push("ROOT")
	return tg.diff(s, other, w)
}

// Find locates the value identified by the path
// Given: <a>x</a>, a.Find({"a"}) returns "x"
func (tg *Tag) Find(path []string) string {
	if len(path) == 0 {
		return tg.Value
	}

	child := path[0]
	for _, ct := range tg.Children {
		if ct.Name == child {
			return ct.Find(path[1:])
		}
	}

	return ""
}
