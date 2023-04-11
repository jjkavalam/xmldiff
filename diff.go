package xmldiff

import (
	"fmt"
	"io"
)

func (tg *Tag) diff(ctxStack *stack, other *Tag, w io.StringWriter) (hasDiff bool) {
	if tg.Name != other.Name {
		must(w.WriteString(fmt.Sprintf("%s TAG: '%s' does not match '%s'\n", bold(ctxStack.String()), red(tg.Name), green(other.Name))))
		hasDiff = true
		return
	}
	ctxStack.push(tg.Name)
	defer ctxStack.pop()
	if len(tg.Children) == 0 && len(other.Children) == 0 {
		if tg.Value != other.Value {
			must(w.WriteString(fmt.Sprintf("%s VALUE: '%s' does not match '%s'\n", bold(ctxStack.String()), red(shorten(tg.Value)), green(shorten(other.Value)))))
			hasDiff = true
		}
		return
	}
	if len(tg.Children) == 0 {
		must(w.WriteString(fmt.Sprintf("%s VALUE: '%s' does not match tag <%s>\n", bold(ctxStack.String()), red(tg.Value), green(other.Name))))
		hasDiff = true
		return
	}
	if len(other.Children) == 0 {
		must(w.WriteString(fmt.Sprintf("%s CHILD_TAGS: <%s>'s child tags are matched by a value '%s'\n", bold(ctxStack.String()), red(tg.Name), green(other.Value))))
		hasDiff = true
		return
	}
	if len(tg.Children) != len(other.Children) {
		must(w.WriteString(fmt.Sprintf("%s CHILD_COUNT: child counts differ %d vs %d\n", bold(ctxStack.String()), len(tg.Children), len(other.Children))))
		hasDiff = true
	}

	var hasDiff2 = false

	// diff children
	hasDiff2 = diffChildren(ctxStack, tg.Children, other.Children, w)

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
