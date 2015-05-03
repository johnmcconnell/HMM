package hmm

import(
	"fmt"
	"bytes"
	"github.com/johnmcconnell/gologspace"
)

type Backward struct {
	tags []Tag
	sequence []string
	filled bool
	trellis *Trellis
	initialState *InitialState
	transition *Transition
	emission *Emission
}

// NewViterb ...
func NewBackward(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *Backward {
	trellis := NewTrellis(tags, len(sequence))
	v := Backward{tags, sequence, false, trellis, i, t, e}
	return &v
}

// String ...
func (v *Backward) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Backward: '%s'\n", v.sequence))
	buffer.WriteString(v.trellis.String())
	return fmt.Sprintf(buffer.String())
}

// ComputeInitialProb ...
func (v *Backward) ComputeInitialProb() float64 {
	return gologspace.LogProb(1.0)
}

// ComputeProb ...
func (v *Backward) ComputeProb(givenTag Tag, index int) float64 {
	value := v.sequence[index + 1]
	pSum := 0.0
	for _, tag := range v.tags {
		nextResult := (*v.trellis)[givenTag][index + 1]
		p := nextResult.Probability
		pT := gologspace.LogProb(v.transition.P(tag, givenTag))
		pE := gologspace.LogProb(v.emission.P(tag, value))

		if pSum == 0.0 {
		  pSum = pE + pT + p
		} else {
		  pSum = gologspace.LogAdd(pSum, pE + pT + p)
		}
	}
	return pSum
}

// Result ...
func (v *Backward) Result(tag Tag, index int) *Result {
	if !v.filled {
		v.FillTrellis()
	}
	return (*v.trellis)[tag][index]
}

// FillTrellis ...
func (v *Backward) FillTrellis() {
	if v.filled {
		return
	}
	l := len(v.sequence) - 1
	for i, _ := range v.sequence {
		v.FillColumn(l - i)
	}
	v.filled = true
}

// FillColumn ...
func (v *Backward) FillColumn(index int) {
	if v.filled {
		return
	}
	for _, tag := range v.tags {
		(*v.trellis)[tag][index] = v.P(tag, index)
	}
}

// P ...
func (v *Backward) P(tag Tag, index int) *Result {
	if index == (len(v.sequence) - 1) {
		p := v.ComputeInitialProb()
		return &Result{"e", p}
	} else {
		p := v.ComputeProb(tag, index)
		return &Result{"e", p}
	}
}
