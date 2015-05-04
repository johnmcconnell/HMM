package hmm

import(
	"fmt"
	"bytes"
	"github.com/johnmcconnell/gologspace"
)

type Forward struct {
	tags []Tag
	sequence []string
	cache *Trellis
	s gologspace.Space
	i *Initial
	t *Transition
	e *Emission
}

// NewViterb ...
func NewForward(tags []Tag, sequence []string, s gologspace.Space, i *Initial, t *Transition, e *Emission) *Forward {
	cache := NewTrellis(tags, len(sequence))
	f := Forward{tags, sequence, cache, s, i, t, e}
	return &f
}

// String ...
func (f *Forward) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Forward: '%s'\n", f.sequence))
	buffer.WriteString(f.cache.FormatString(f.s))
	return fmt.Sprintf(buffer.String())
}

// CompInitialP ...
func (f *Forward) CompInitialP(pI, pE float64) float64 {
	s := f.s
	return s.Mul(pI, pE)
}

// CompTransP ...
func (f *Forward) CompTransP(pT, pP float64) float64 {
	s := f.s
	return s.Mul(pT, pP)
}

// CompProb ...
func (f *Forward) CompP(tag Tag, index int, pE float64) float64 {
	pSum := 0.0
	for _, givenTag := range f.tags {
		prevResult := (*f.cache)[givenTag][index - 1]
		prevP := prevResult.Prob
		pT := f.t.P(tag, givenTag)
		p := f.CompTransP(pT, prevP)
		if (pSum == 0.0) {
			pSum = p
		} else {
			pSum = f.s.Add(pSum, p)
		}
	}
	return f.s.Mul(pE, pSum)
}

// FillTrellis ...
func (f *Forward) FillTrellis() {
	for i, _ := range f.sequence {
		f.FillColumn(i)
	}
}

// FillColumn ...
func (f *Forward) FillColumn(index int) {
	for _, tag := range f.tags {
	  f.FillValue(tag, index)
	}
}

// FillValue ...
func (f *Forward) FillValue(tag Tag, i int) {
	v := f.sequence[i]
	pE := f.e.P(tag, v)
	if i == 0 {
		pI := f.i.P(tag)
		p := f.CompInitialP(pI, pE)
		(*f.cache)[tag][i] = &Result{"e", p}
	} else {
		p := f.CompP(tag, i, pE)
		(*f.cache)[tag][i] = &Result{"e", p}
	}
}

// P ...
func (f *Forward) R(tag Tag, index int) *Result {
	return (*f.cache)[tag][index]
}
