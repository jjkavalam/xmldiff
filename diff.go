package xmldiff

import (
	"fmt"
	"io"
	"strings"
)

func (tg *Tag) diff(ctxStack *Stack, other *Tag, w io.StringWriter) error {
	ctx := strings.Join(*ctxStack, ">")
	if tg.Name != other.Name {
		_, err := w.WriteString(fmt.Sprintf("%s TAG: '%s' is matched by '%s'\n", ctx, tg.Name, other.Name))
		if err != nil {
			return err
		}
		return nil
	}
	ctxStack.Push(tg.Name)
	defer ctxStack.Pop()
	ctx = strings.Join(*ctxStack, ">")
	if len(tg.Children) == 0 && len(other.Children) == 0 {
		if tg.Value != other.Value {
			_, err := w.WriteString(fmt.Sprintf("%s VALUE: '%s' is matched by '%s'\n", ctx, tg.Value, other.Value))
			if err != nil {
				return err
			}
		}
		return nil
	}
	if len(tg.Children) == 0 {
		_, err := w.WriteString(fmt.Sprintf("%s VALUE: '%s' is matched by a tag <%s>\n", ctx, tg.Value, other.Name))
		if err != nil {
			return err
		}
		return nil
	}
	if len(other.Children) == 0 {
		_, err := w.WriteString(fmt.Sprintf("%s CHILD_TAGS: <%s>'s child tags are matched by a value '%s'\n", ctx, tg.Name, other.Value))
		if err != nil {
			return err
		}
		return nil
	}
	if len(tg.Children) != len(other.Children) {
		_, err := w.WriteString(fmt.Sprintf("%s CHILD_COUNT: child counts differ %d vs %d\n", ctx, len(tg.Children), len(other.Children)))
		if err != nil {
			return err
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
				err := ctg.diff(ctxStack, oTg, w)
				if err != nil {
					return err
				}
				// since match is made; remove both from unmatched sets
				delete(unmatchedRight, k)
				delete(unmatchedLeft, i)
				break
			}
		}
	}

	for k := range unmatchedLeft {
		_, err := w.WriteString(fmt.Sprintf("%s REMOVED_TAG: %s\n", ctx, tg.Children[k].Name))
		if err != nil {
			return err
		}
	}
	for k := range unmatchedRight {
		_, err := w.WriteString(fmt.Sprintf("%s ADDED_TAG: %s\n", ctx, other.Children[k].Name))
		if err != nil {
			return err
		}
	}
	return nil
}
