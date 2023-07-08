package main

import (
	"fmt"
)

func main() {
	str := "hello 世界上"
	fmt.Println(len(str))
	bytes := []byte(str)
	fmt.Println(len(bytes))
}
