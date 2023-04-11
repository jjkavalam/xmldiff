package xmldiff

import (
	"fmt"
	"strings"
)

type stack []string

func (s *stack) push(item string) {
	*s = append(*s, item)
}

func (s *stack) pop() string {
	if len(*s) == 0 {
		return ""
	}
	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item
}

func (s *stack) String() string {
	return fmt.Sprintf("[%s]\n", strings.Join(*s, ">"))
}

func newStack() *stack {
	s := stack([]string{})
	return &s
}
