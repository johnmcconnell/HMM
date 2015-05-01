package main

import (
	"fmt"
	"github.com/johnmcconnell/hmm/src/trellis"
)

func main() {
	tags := []trellis.Tag{"A", "B", "C"}
	size := 10
	initialState := trellis.InitialState{}
	transition := trellis.Transition{}
	emission := trellis.Emission{}
	t := trellis.New(tags, size, initialState, transition, emission)
	fmt.Printf("%s\n", t.Test())
}
