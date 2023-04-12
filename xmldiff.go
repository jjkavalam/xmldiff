package xmldiff

import (
	"fmt"
	"io"
)

func Parse(xmlData string) (*Tag, error) {
	p := newParser(xmlData)
	return p.parseTag()
}

type Tag struct {
	Name     string
	Children []*Tag
	// Value is valid only if Children is empty
	Value string
}

func (tg *Tag) StringChildrenOnly(w io.StringWriter) {
	if len(tg.Children) == 0 {
		return
	}
	for _, ctg := range tg.Children {
		ctg.String(w)
	}
}

func (tg *Tag) String(w io.StringWriter) {
	must(w.WriteString(fmt.Sprintf("<%s>", tg.Name)))
	if len(tg.Children) == 0 {
		must(w.WriteString(tg.Value))
	} else {
		tg.StringChildrenOnly(w)
	}
	must(w.WriteString(fmt.Sprintf("</%s>", tg.Name)))
}

// Diff compares this tag with another.
// It performs a tree traversal on both trees simultaneously and returns a list of differences between the trees.
func (tg *Tag) Diff(other *Tag, w io.StringWriter) bool {
	s := newStack()
	s.push("ROOT")
	return tg.diff(s, other, w)
}

// Find locates the value identified by the path
// Given: <a>x</a>, a.Find({"a"}) returns *Tag<a>
func (tg *Tag) Find(path []string) *Tag {
	if len(path) == 0 {
		return tg
	}

	child := path[0]
	for _, ct := range tg.Children {
		if ct.Name == child {
			return ct.Find(path[1:])
		}
	}

	return nil
}
