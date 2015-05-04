package hmm

import(
	"fmt"
	"bytes"
	"github.com/johnmcconnell/gologspace"
)

type Backward struct {
	tags []Tag
	sequence []string
	cache *Trellis
	s gologspace.Space
	i *Initial
	t *Transition
	e *Emission
}

// NewBackward ...
func NewBackward(tags []Tag, sequence []string, s gologspace.Space, i *Initial, t *Transition, e *Emission) *Backward {
	cache := NewTrellis(tags, len(sequence))
	b := Backward{tags, sequence, cache, s, i, t, e}
	return &b
}

// String ...
func (b *Backward) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Backward: '%s'\n", b.sequence))
	buffer.WriteString(b.cache.FormatString(b.s))
	return fmt.Sprintf(buffer.String())
}

// ComputeInitialProb ...
func (b *Backward) CompInitialP() float64 {
	return b.s.Enter(1.0)
}

// ComputeTransitionProb
func (b *Backward) CompTransP(pE, pT, pP float64) float64 {
	s := b.s
	return s.Mul(s.Mul(pE, pT), pP)
}

// ComputeProb ...
func (b *Backward) CompP(givenTag Tag, index int) float64 {
	value := b.sequence[index + 1]
	pSum := 0.0
	s := b.s
	for _, tag := range b.tags {
		nextR := (*b.cache)[givenTag][index + 1]
		pP := nextR.Prob
		pT := b.t.P(tag, givenTag)
		pE := b.e.P(tag, value)
		p := b.CompTransP(pE, pT, pP)
		if (pSum == 0.0) {
			pSum = p
		} else {
			pSum = s.Add(pSum, p)
		}
	}
	return pSum
}

// FillTrellis ...
func (b *Backward) FillTrellis() {
	l := len(b.sequence) - 1
	for i, _ := range b.sequence {
		rI := l - i
		b.FillColumn(rI)
	}
}

// FillColumn ...
func (b *Backward) FillColumn(index int) {
	for _, tag := range b.tags {
		b.FillValue(tag, index)
	}
}

// FillValue ...
func (b *Backward) FillValue(tag Tag, index int) {
	if index == (len(b.sequence) - 1) {
		p := b.CompInitialP()
		(*b.cache)[tag][index] = &Result{"e", p}
	} else {
		p := b.CompP(tag, index)
		(*b.cache)[tag][index] = &Result{"e", p}
	}
}

// Result ...
func (b *Backward) R(tag Tag, index int) *Result {
	return (*b.cache)[tag][index]
}
