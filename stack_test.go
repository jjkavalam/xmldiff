package xmldiff_test

import (
	"github.com/jjkavalam/xmldiff"
	"strings"
	"testing"
)

func TestStack(t *testing.T) {
	s := xmldiff.NewStack()

	s.Push("a")
	s.Push("b")
	s.Push("c")
	s.Pop()
	s.Push("d")

	if strings.Join(*s, ":") != "a:b:d" {
		t.Errorf("want 'a:b:d'; got %s", *s)
	}
}
