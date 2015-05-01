package main

import (
	"fmt"
	"github.com/johnmcconnell/trellis"
)

func main() {
	tags := []string{"A", "B", "C"}
	size := 10
	t := trellis.New(tags, size)
	fmt.Printf(t)
}
