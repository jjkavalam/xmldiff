package xmldiff_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"

	"github.com/jjkavalam/xmldiff"
)

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

	err = t1.String(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n---")

	t2, err := xmldiff.Parse(xmlData2)
	if err != nil {
		log.Fatal(err)
	}

	err = t2.String(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n---")

	var outBuf bytes.Buffer
	err = t1.Diff(t2, &outBuf)
	if err != nil {
		log.Fatal(err)
	}

	t.Log(outBuf.String())

	expected := `ROOT>x CHILD_COUNT: child counts differ 2 vs 3
ROOT>x>c>e VALUE: 'hello 
world' is matched by 'g'
ROOT>x REMOVED_TAG: b
ROOT>x ADDED_TAG: d
ROOT>x ADDED_TAG: d
`

	assert.Equal(t, expected, outBuf.String())

}
