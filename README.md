# xmldiff

See `xmldiff_test.go` for usage.

Key structs, functions and methods:

```go
type Tag struct {
	Name     string
	Children []*Tag
	// Value is valid only if Children is empty
	Value string
}

func Parse(xmlData string) (*Tag, error)
func (tg *Tag) String(w io.StringWriter) error
func (tg *Tag) Diff(other *Tag, w io.StringWriter) error
```