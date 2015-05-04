package hmm

import(
	"github.com/johnmcconnell/gologspace"
)

type Transition map[Tag]map[Tag]float64

func (t *Transition) P(tag, givenTag Tag) float64 {
	return (*t)[givenTag][tag]
}

func NewTransition(tags []Tag) Transition {
	t := make(Transition)
	for _, tag := range tags {
		t[tag] = make(map[Tag]float64)
	}
	return t
}

func UniformT(tags []Tag, s gologspace.Space) Transition {
	t := make(Transition)
	length := float64(len(tags))
	size := s.Enter(length)
  one := s.Enter(1.0)
	for _, tag1 := range tags {
		t[tag1] = make(map[Tag]float64)
		for _, tag2 := range tags {
			t[tag1][tag2] = one / size
		}
	}
	return t
}

func UniformT2(tags TCache, s gologspace.Space) Transition {
	t := make(Transition)
  one := s.Enter(1.0)
	for tag1, nTags := range tags {
		t[tag1] = make(map[Tag]float64)
	  length := float64(len(nTags))
		size := s.Enter(length)
		for tag2, _ := range nTags {
			t[tag1][tag2] = one / size
		}
	}
	return t
}
