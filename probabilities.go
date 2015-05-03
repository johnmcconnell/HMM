package hmm

type InitialState map[Tag]float64
type Transition map[Tag]map[Tag]float64
type Emission map[Tag]map[string]float64

func (i *InitialState) P(tag Tag) float64 {
	return (*i)[tag]
}

func (t *Transition) P(tag Tag, givenTag Tag) float64 {
	return (*t)[givenTag][tag]
}

func (e *Emission) P(tag Tag, value string) float64 {
	return (*e)[tag][value]
}

func NewEmission(tags []Tag) Emission {
	e := make(Emission)
	for _, tag := range tags {
		e[tag] = make(map[string]float64)
	}
	return e
}

func NewTransition(tags []Tag) Transition {
	t := make(Transition)
	for _, tag := range tags {
		t[tag] = make(map[Tag]float64)
	}
	return t
}

func UniformI(tags []Tag) InitialState {
	iS := make(InitialState)
	l := len(tags)
	for _, tag := range tags {
		iS[tag] = 1.0 / float64(l)
	}
	return iS
}

func UniformT(tags []Tag) Transition {
	t := make(Transition)
	l := len(tags)
	for _, tag1 := range tags {
		t[tag1] = make(map[Tag]float64)
		for _, tag2 := range tags {
			t[tag1][tag2] = 1.0 / float64(l)
		}
	}
	return t
}

func UniformE2(tags []Tag, words map[string]bool) Emission {
	e := make(Emission)
	for _, tag := range tags {
		e[tag] = make(map[string]float64)
		l := len(words)
		for word, _ := range words {
			e[tag][word] = 1.0 / float64(l)
		}
	}
	return e
}

func UniformE(possibleEmissions map[Tag]map[string]bool) Emission {
	e := make(Emission)
	for tag, words := range possibleEmissions {
		e[tag] = make(map[string]float64)
		l := len(words)
		for word, _ := range words {
			e[tag][word] = 1.0 / float64(l)
		}
	}
	return e
}
