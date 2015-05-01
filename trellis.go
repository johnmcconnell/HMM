package main

import(
	"fmt"
)

type Result struct {
	previousTag Tag
	probability float64
}

type Trellis map[Tag][]*Result

// String ...
func (t *Trellis) String() string {
	return fmt.Sprintf("Hello")
}

// Set ...
func (r *Result) Set(r2 Result) {
	*r = r2
}

func MakeResult(tag Tag, probability float64) Result {
	return Result{tag, probability}
}

// New ...
func New(tags []Tag, size int) *Trellis {
	t := make(Trellis)
	for _, tag := range tags {
		t[tag] = make([]*Result, size)
		for i,_ := range t[tag] {
			t[tag][i] = new(Result)
		}
	}
	return &t
}

