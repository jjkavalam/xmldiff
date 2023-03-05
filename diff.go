package xmldiff

import (
	"fmt"
	"io"
	"strings"
)

func (tg *Tag) diff(ctxStack *Stack, other *Tag, w io.StringWriter) (hasDiff bool, err error) {
	ctx := strings.Join(*ctxStack, ">")
	if tg.Name != other.Name {
		_, err = w.WriteString(fmt.Sprintf("%s TAG: '%s' is matched by '%s'\n", Bold(ctx), Red(tg.Name), Green(other.Name)))
		if err != nil {
			return
		}
		hasDiff = true
		return
	}
	ctxStack.Push(tg.Name)
	defer ctxStack.Pop()
	ctx = strings.Join(*ctxStack, ">")
	if len(tg.Children) == 0 && len(other.Children) == 0 {
		if tg.Value != other.Value {
			_, err = w.WriteString(fmt.Sprintf("%s VALUE: '%s' is matched by '%s'\n", Bold(ctx), Red(shorten(tg.Value)), Green(shorten(other.Value))))
			if err != nil {
				return
			}
			hasDiff = true
		}
		return
	}
	if len(tg.Children) == 0 {
		_, err = w.WriteString(fmt.Sprintf("%s VALUE: '%s' is matched by a tag <%s>\n", Bold(ctx), Red(tg.Value), Green(other.Name)))
		if err != nil {
			return
		}
		hasDiff = true
		return
	}
	if len(other.Children) == 0 {
		_, err = w.WriteString(fmt.Sprintf("%s CHILD_TAGS: <%s>'s child tags are matched by a value '%s'\n", Bold(ctx), Red(tg.Name), Green(other.Value)))
		if err != nil {
			return
		}
		hasDiff = true
		return
	}
	if len(tg.Children) != len(other.Children) {
		_, err = w.WriteString(fmt.Sprintf("%s CHILD_COUNT: child counts differ %d vs %d\n", Bold(ctx), len(tg.Children), len(other.Children)))
		hasDiff = true
		if err != nil {
			return
		}
	}

	// diff children
	unmatchedLeft := map[int]bool{}
	for i := range tg.Children {
		unmatchedLeft[i] = true
	}

	unmatchedRight := map[int]bool{}
	for i := range other.Children {
		unmatchedRight[i] = true
	}

	for i, ctg := range tg.Children {
		// try to match with an unmatchedRight tag
		for k := range unmatchedRight {
			oTg := other.Children[k]
			if ctg.Name == oTg.Name {
				hasDiff2, err := ctg.diff(ctxStack, oTg, w)
				if err != nil {
					return false, err
				}
				if hasDiff2 {
					hasDiff = true
				}
				// since match is made; remove both from unmatched sets
				delete(unmatchedRight, k)
				delete(unmatchedLeft, i)
				break
			}
		}
	}

	for k := range unmatchedLeft {
		_, err = w.WriteString(fmt.Sprintf("%s REMOVED_TAG: %s\n", Bold(ctx), Red(tg.Children[k].Name)))
		if err != nil {
			return
		}
		hasDiff = true
	}
	for k := range unmatchedRight {
		_, err = w.WriteString(fmt.Sprintf("%s ADDED_TAG: %s\n", Bold(ctx), Green(other.Children[k].Name)))
		if err != nil {
			return
		}
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
