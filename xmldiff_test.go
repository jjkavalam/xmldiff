package xmldiff_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jjkavalam/xmldiff"
)

func TestFind(t *testing.T) {

	data := `
<?xml version="1.0" encoding="utf-8"?>
<a>
<c>hello</c>
</a> `
	tag, err := xmldiff.Parse(data)

	if err != nil {
		t.Fatal(err)
	}

	v := tag.Find([]string{"c"})

	if v.Value != "hello" {
		t.Fatalf("expected '%s', got '%s'", "hello", v)
	}

}

func TestDiff(t *testing.T) {
	t.Setenv("NO_COLOR", "true")

	xmlData1 := `<x>

<c><e>hello 
world</e></c>
<b>hello</b>
</x> `

	xmlData2 := `<x>

<d></d>
<c><e>g</e>
</c><d>ok</d></x>`

	t1, err := xmldiff.Parse(xmlData1)
	if err != nil {
		log.Fatal(err)
	}

	t1.String(os.Stdout)

	fmt.Println("\n---")

	t2, err := xmldiff.Parse(xmlData2)
	if err != nil {
		log.Fatal(err)
	}

	t2.String(os.Stdout)

	fmt.Println("\n---")

	var outBuf bytes.Buffer
	hasDiff := t1.Diff(t2, &outBuf)

	if !hasDiff {
		t.Fatal("expected hasDiff to be true")
	}

	t.Log(outBuf.String())

	expected := `[ROOT>x]
 CHILD_COUNT: child counts differ 2 vs 3
[ROOT>x]
 ADDED_TAG: d
[ROOT>x>c>e]
 VALUE: 'hello 
world' does not match 'g'
[ROOT>x]
 REMOVED_TAG: b
[ROOT>x]
 ADDED_TAG: d
`

	if expected != outBuf.String() {
		t.Fatalf("want '%s', got '%s'", expected, outBuf.String())
	}

}
