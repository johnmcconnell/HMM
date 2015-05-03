package hmm

import(
	"fmt"
	"bytes"
)

type BackwardLog struct {
	tags []Tag
	sequence []string
	filled bool
	trellis *Trellis
	initialState *InitialState
	transition *Transition
	emission *Emission
}

// NewViterb ...
func NewBackwardLog(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *BackwardLog {
	trellis := NewTrellis(tags, len(sequence))
	v := BackwardLog{tags, sequence, false, trellis, i, t, e}
	return &v
}

// String ...
func (v *BackwardLog) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("BackwardLog: '%s'\n", v.sequence))
	buffer.WriteString(v.trellis.String())
	return fmt.Sprintf(buffer.String())
}

// ComputeInitialProb ...
func (v *BackwardLog) ComputeInitialProb() float64 {
	return 1.0
}

// ComputeProb ...
func (v *BackwardLog) ComputeProb(givenTag Tag, index int) float64 {
	value := v.sequence[index + 1]
	pSum := 0.0
	for _, tag := range v.tags {
		nextResult := (*v.trellis)[givenTag][index + 1]
		p := nextResult.Probability
		pT := v.transition.P(tag, givenTag)
		pE := v.emission.P(tag, value)
		pSum += pE * pT * p
	}
	return pSum
}

// Result ...
func (v *BackwardLog) Result(tag Tag, index int) *Result {
	if !v.filled {
		v.FillTrellis()
	}
	return (*v.trellis)[tag][index]
}

// FillTrellis ...
func (v *BackwardLog) FillTrellis() {
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
func (v *BackwardLog) FillColumn(index int) {
	if v.filled {
		return
	}
	for _, tag := range v.tags {
		(*v.trellis)[tag][index] = v.P(tag, index)
	}
}

// P ...
func (v *BackwardLog) P(tag Tag, index int) *Result {
	if index == (len(v.sequence) - 1) {
		p := v.ComputeInitialProb()
		return &Result{"e", p}
	} else {
		p := v.ComputeProb(tag, index)
		return &Result{"e", p}
	}
}
