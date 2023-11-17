package main

import (
	"db/file"
	"fmt"
)

func main() {
	g := []string{}
	g = append(g, "dasdasdas")
	fmt.Println(file.CodeingToRow(g))
}
