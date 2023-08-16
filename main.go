package main

import (
	"encoding/json"
	"fmt"
	"github.com/iancoleman/orderedmap"
)

func main() {
	m := orderedmap.New()
	s := "{\n  \"IDs\": [\n        7236290603911250220\n    ] \n}"

	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}
