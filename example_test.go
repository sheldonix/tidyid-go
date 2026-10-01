package tidyid_test

import (
	"fmt"
	"log"

	"github.com/sheldonix/tidyid-go/v2"
)

func ExampleGenerate() {
	id, err := tidyid.Generate(16, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tidyid.IsValidIDOfLength(id, 16, false))
	// Output: true
}

func ExampleGenerate_uppercase() {
	id, err := tidyid.Generate(10, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tidyid.IsValidIDOfLength(id, 10, true))
	// Output: true
}
