package trellis

import(
	"fmt"
)

type TrellisResult struct {
	previousTag string
	probability float64
}

type InitialState map[string]float64
type Transition map[Tag]map[Tag]float64
type Emission map[Tag]map[string]float64

type Tag string

type Trellis struct {
	cache map[Tag][]*TrellisResult
	initialState InitialState
	transition Transition
	emission Emission
}

// String ...
func (t *Trellis) String() string {
	return fmt.Sprintf("Hello")
}

func (t *Trellis) Test() string {
	return "blah"
}

// New ...
func New(tags []Tag, size int, initialState InitialState,
transition Transition, emission Emission) *Trellis {
	cache := make(map[Tag][]*TrellisResult)
	for _, tag := range tags {
		cache[tag] = make([]*TrellisResult, size)
	}
	return 2
	return &Trellis{cache, initialState, transition, emission}
}
