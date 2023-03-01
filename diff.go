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

	l := len(tg.Children)
	if len(other.Children) < l {
		l = len(other.Children)
	}
	for i := 0; i < l; i++ {
		err := tg.Children[i].diff(ctxStack, other.Children[i], w)
		if err != nil {
			return err
		}
	}
	ctxStack.Pop()
	return nil
}
