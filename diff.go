package xmldiff

import (
	"fmt"
	"io"
)

func (tg *Tag) diff(ctxStack *Stack, other *Tag, w io.StringWriter) (hasDiff bool, err error) {
	if tg.Name != other.Name {
		_, err = w.WriteString(fmt.Sprintf("%s TAG: '%s' does not match '%s'\n", Bold(ctxStack.String()), Red(tg.Name), Green(other.Name)))
		if err != nil {
			return
		}
		hasDiff = true
		return
	}
	ctxStack.Push(tg.Name)
	defer ctxStack.Pop()
	if len(tg.Children) == 0 && len(other.Children) == 0 {
		if tg.Value != other.Value {
			_, err = w.WriteString(fmt.Sprintf("%s VALUE: '%s' does not match '%s'\n", Bold(ctxStack.String()), Red(shorten(tg.Value)), Green(shorten(other.Value))))
			if err != nil {
				return
			}
			hasDiff = true
		}
		return
	}
	if len(tg.Children) == 0 {
		_, err = w.WriteString(fmt.Sprintf("%s VALUE: '%s' does not match tag <%s>\n", Bold(ctxStack.String()), Red(tg.Value), Green(other.Name)))
		if err != nil {
			return
		}
		hasDiff = true
		return
	}
	if len(other.Children) == 0 {
		_, err = w.WriteString(fmt.Sprintf("%s CHILD_TAGS: <%s>'s child tags are matched by a value '%s'\n", Bold(ctxStack.String()), Red(tg.Name), Green(other.Value)))
		if err != nil {
			return
		}
		hasDiff = true
		return
	}
	if len(tg.Children) != len(other.Children) {
		_, err = w.WriteString(fmt.Sprintf("%s CHILD_COUNT: child counts differ %d vs %d\n", Bold(ctxStack.String()), len(tg.Children), len(other.Children)))
		hasDiff = true
		if err != nil {
			return
		}
	}

	var hasDiff2 = false

	// diff children
	hasDiff2, err = diffChildren(ctxStack, tg.Children, other.Children, w)
	if err != nil {
		return
	}

	if hasDiff2 {
		hasDiff = true
	}

	return
}

func shorten(s string) string {
	if len(s) < OptionShortenValueDiffs {
		return s
	}
	return s[:OptionShortenValueDiffs] + "..."
}
