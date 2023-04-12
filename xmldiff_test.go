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
<e></e>
<f><g></g></f>
</x> `

	xmlData2 := `<x>

<d></d>
<c><e>g</e>
</c><d>ok</d>
<e><g></g></e>
<f>xyz</f>
</x>
`

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

	expected := `[ROOT>x]
 CHILD_COUNT: child counts differ 4 vs 5
[ROOT>x]
 ADDED_TAG: d (found at position 0)
[ROOT>x>c>e]
 VALUE: 'hello 
world' does not match 'g'
[ROOT>x]
 REMOVED_TAG: b (expected at position 1)
[ROOT>x]
 ADDED_TAG: d (found at position 2)
[ROOT>x>e]
 VALUE: '' does not match '<g></g>'
[ROOT>x>f]
 CHILD_TAGS: '<g></g>' does not match 'xyz'
`

	actual := outBuf.String()
	if expected != actual {
		t.Fatalf("want '%s', got '%s'", expected, actual)
	}

}
