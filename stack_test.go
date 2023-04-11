package xmldiff

import (
	"strings"
	"testing"
)

func TestStack(t *testing.T) {
	s := newStack()

	s.push("a")
	s.push("b")
	s.push("c")
	s.pop()
	s.push("d")

	if strings.Join(*s, ":") != "a:b:d" {
		t.Errorf("want 'a:b:d'; got %s", *s)
	}
}
