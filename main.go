package main

import (
	"fmt"
	"strings"
)

var users []string = []string{
	"ZHANGYIMAANG",
	"liangrubo",
}

type Employee struct {
	Name       string
	Department string
	IsMember   bool
}

func main() {
	uri := ""
	res := strings.SplitN(uri, "/", 2)
	fmt.Println(res)
}
