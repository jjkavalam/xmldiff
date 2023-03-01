package xmldiff

type Stack []string

func (s *Stack) Push(item string) {
	*s = append(*s, item)
}

func (s *Stack) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item
}

func NewStack() *Stack {
	s := Stack([]string{})
	return &s
}
