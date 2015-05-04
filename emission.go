package hmm

import(
	"github.com/johnmcconnell/gologspace"
)

type Emission map[Tag]map[string]float64

func (e *Emission) P(tag Tag, val string) float64 {
	return (*e)[tag][val]
}

func NewEmission(tags []Tag) Emission {
	e := make(Emission)
	for _, tag := range tags {
		e[tag] = make(map[string]float64)
	}
	return e
}

func UniformE(emissions ECache, words []string, s gologspace.Space) Emission {
	e := make(Emission)
  one := s.Enter(1.0)
	zero := s.Enter(0.0)
	for tag, foundWs := range emissions {
		e[tag] = make(map[string]float64)
		// Zero initial values
		for _, word := range words {
			e[tag][word] = zero
		}
	  length := float64(len(words))
		size := s.Enter(length)
		for word, _ := range foundWs {
			e[tag][word] = one / size
		}
	}
	return e
}

func UniformE2(tags []Tag, words []string, s gologspace.Space) Emission {
	e := make(Emission)
  one := s.Enter(1.0)
	for _, tag := range tags {
		e[tag] = make(map[string]float64)
	  length := float64(len(words))
		size := s.Enter(length)
		for _, word := range words {
			e[tag][word] = one / size
		}
	}
	return e
}
