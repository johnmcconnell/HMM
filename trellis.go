package hmm

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

// New ...
func New(tags []Tag, size int,
initialState InitialState, transition Transition, emission Emission) *Trellis {
	cache := make(map[Tag][]*TrellisResult)
	for _, tag := range tags {
		cache[tag] = make([]*TrellisResult, size)
	}
	return &Trellis{cache, initialState, transition, emission}
}
