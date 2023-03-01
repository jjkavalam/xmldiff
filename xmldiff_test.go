package xmldiff_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jjkavalam/xmldiff"
)

func TestMain(t *testing.T) {

	xmlData1 := `<x>

<b>hello</b>
<c></c>
</x> `

	xmlData2 := `<a><c></c></a>`

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

	err = t1.Diff(t2, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
