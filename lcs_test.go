package xmldiff

import (
	"testing"
)

func TestLongestCommonSubsequence(t *testing.T) {
	a := []*Tag{
		newTag("<a></a>"),
		newTag("<b></b>"),
		newTag("<c></c>"),
		newTag("<d></d>"),
	}

	b := []*Tag{
		newTag("<a></a>"),
		newTag("<c></c>"),
		newTag("<e></e>"),
	}
	r := longestCommonSubsequence[*Tag](a, b, func(a, b *Tag) bool {
		return a.Name == b.Name
	})
	if len(r) != 2 {
		t.Fatal("wrong length")
	}
	if r[0].Name != "a" || r[1].Name != "c" {
		t.Fatal("wrong subsequence")
	}
}

func newTag(xmlData string) *Tag {
	tag, err := Parse(xmlData)
	if err != nil {
		panic(err)
	}
	return tag
}
