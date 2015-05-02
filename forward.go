package hmm

import(
	"fmt"
	"bytes"
)

type Forward struct {
	tags []Tag
	sequence string
	filled bool
	trellis *Trellis
	initialState *InitialState
	transition *Transition
	emission *Emission
}

// NewViterb ...
func NewForward(tags []Tag, sequence string, i *InitialState, t *Transition, e *Emission) *Forward {
	trellis := NewTrellis(tags, len(sequence))
	v := Forward{tags, sequence, false, trellis, i, t, e}
	return &v
}

func (v *Forward) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Forward: '%s'\n", v.sequence))
	buffer.WriteString(v.trellis.String())
	return fmt.Sprintf(buffer.String())
}

// Result ...
func (v *Forward) Result(tag Tag, index int) *Result {
	if !v.filled {
		v.FillTrellis()
	}
	return (*v.trellis)[tag][index]
}

// FillTrellis ...
func (v *Forward) FillTrellis() {
	if v.filled {
		return
	}
	for i, _ := range v.sequence {
		v.FillColumn(i)
	}
	v.filled = true
}

// FillColumn ...
func (v *Forward) FillColumn(index int) {
	if v.filled {
		return
	}
	for _, tag := range v.tags {
		(*v.trellis)[tag][index] = v.P(tag, index)
	}
}

// P ...
func (v *Forward) P(tag Tag, index int) *Result {
	if index == 0 {
		value := v.sequence[index]
		pI := v.initialState.P(tag)
		pE := v.emission.P(tag, value)
		return &Result{"e", pI * pE}
	} else {
	  return v.SumP(tag, index)
	}
}

// SumP ...
func (v *Forward) SumP(tag Tag, index int) *Result {
  value := v.sequence[index]
	sumResult := &Result{"e", 0.0}
	for _, givenTag := range v.tags {
		prevResult := (*v.trellis)[givenTag][index - 1]
		p := prevResult.Probability
		pT := v.transition.P(tag, givenTag)
		sumResult.Probability += pT * p
	}
  pE := v.emission.P(tag, value)
	sumResult.Probability = sumResult.Probability * pE
	return sumResult
}
