package xmldiff

import (
	"fmt"
	"io"
)

func diffChildren(ctxStack *stack, this []*Tag, that []*Tag, w io.StringWriter) (hasDiff bool) {
	// two sections that have the same name identify points in the script that needs to match
	compareFn := func(a, b *Tag) bool {
		return a.Name == b.Name
	}

	// find the longest common subsequence of matching sections
	common := longestCommonSubsequence[*Tag](this, that, compareFn)

	// keep a cursor for each side; advance each to the next match in the common sequence;
	// anything skipped on this was deleted; while anything skipped on that was added
	i := 0
	j := 0

	var hasDiff2 = false

	removedTag := func(i int) string {
		t := this[i]
		return fmt.Sprintf("%s REMOVED_TAG: %s (expected at position %d)\n", bold(ctxStack.String()), red(t.Name), i)
	}

	addedTag := func(j int) string {
		t := that[j]
		return fmt.Sprintf("%s ADDED_TAG: %s (found at position %d)\n", bold(ctxStack.String()), green(t.Name), j)
	}

	for k := range common {

		for ; i < len(this); i++ {
			if compareFn(this[i], common[k]) {
				break
			}
			must(w.WriteString(removedTag(i)))
			hasDiff = true
		}

		for ; j < len(that); j++ {
			if compareFn(that[j], common[k]) {
				break
			}
			must(w.WriteString(addedTag(j)))
			hasDiff = true
		}

		if i < len(this) && j < len(that) {
			// compare
			hasDiff2 = this[i].diff(ctxStack, that[j], w)
			if hasDiff2 {
				// the difference would have been already written into w
				hasDiff = true
			} else {
				// there is no difference; print nothing
			}
		}

		i++
		j++

	}

	for ; i < len(this); i++ {
		must(w.WriteString(removedTag(i)))
		hasDiff = true
	}

	for ; j < len(that); j++ {
		must(w.WriteString(addedTag(j)))
		hasDiff = true
	}

	return hasDiff
}
