package hmm

import(
	"log"
)

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
	sum := 0.0
	for _, tag := range tags {
		iS[tag] = 1.0 / float64(l)
		sum += iS[tag]
	}
	log.Printf("S total %v\n", sum)
	return iS
}

func UniformT(tags []Tag) Transition {
	t := make(Transition)
	l := float64(len(tags))
	for _, tag1 := range tags {
		t[tag1] = make(map[Tag]float64)
		sum := 0.0
		for _, tag2 := range tags {
			t[tag1][tag2] = 1.0 / l
			sum += t[tag1][tag2]
		}
		log.Printf("T(%s) total %v\n", tag1, sum)
	}
	return t
}

func UniformE2(tags []Tag, words map[string]bool) Emission {
	e := make(Emission)
	for _, tag := range tags {
		e[tag] = make(map[string]float64)
		l := float64(len(words))
		sum := 0.0
		for word, _ := range words {
			e[tag][word] = 1.0 / l
			sum += e[tag][word]
		}
		log.Printf("E(%s) total %v\n", tag, sum)
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
