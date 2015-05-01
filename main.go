package main

import (
	"fmt"
)

func main() {
	tags := []Tag{"A", "B", "C"}
	size := 10
	// initialState := InitialState{}
	// transition := Transition{}
	// emission := Emission{}
	t := New(tags, size)
	// r := MakeResult(tags[0], 0.312)
	fmt.Printf("%s\n", t)
}
