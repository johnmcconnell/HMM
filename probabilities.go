package hmm

type InitialState map[Tag]float64
type Transition map[Tag]map[Tag]float64
type Emission map[Tag]map[uint8]float64

func (i *InitialState) P(tag Tag) float64 {
	return (*i)[tag]
}

func (t *Transition) P(tag Tag, givenTag Tag) float64 {
	return (*t)[givenTag][tag]
}

func (e *Emission) P(tag Tag, value uint8) float64 {
	return (*e)[tag][value]
}
