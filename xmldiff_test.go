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

	xmlData1 := `<x>

<c><e>f</e></c>
<b>hello</b>
</x> `

	xmlData2 := `<x>
<c><e>g</e>
</c>
<d></d>
</x>`

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

	expected := `ROOT>x>c>e VALUE: 'f' is matched by 'g'
ROOT>x TAG: 'b' is matched by 'd'
`

	assert.Equal(t, expected, outBuf.String())

}
