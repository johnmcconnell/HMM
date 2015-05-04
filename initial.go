package hmm

import(
	"github.com/johnmcconnell/gologspace"
)

type Initial map[Tag]float64

func (i *Initial) P(tag Tag) float64 {
	return (*i)[tag]
}

func UniformI(tags []Tag, s gologspace.Space) Initial {
	iS := make(Initial)
	length := float64(len(tags))
	size := s.Enter(length)
  one := s.Enter(1.0)
	for _, tag := range tags {
		iS[tag] = one / size
	}
	return iS
}
