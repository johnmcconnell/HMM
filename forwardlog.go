package hmm

import(
	"fmt"
	"bytes"
)

type ForwardLog struct {
	tags []Tag
	sequence []string
	filled bool
	trellis *Trellis
	initialState *InitialState
	transition *Transition
	emission *Emission
}

// NewViterb ...
func NewForwardLog(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *ForwardLog {
	trellis := NewTrellis(tags, len(sequence))
	v := ForwardLog{tags, sequence, false, trellis, i, t, e}
	return &v
}

// String ...
func (v *ForwardLog) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("ForwardLog: '%s'\n", v.sequence))
	buffer.WriteString(v.trellis.String())
	return fmt.Sprintf(buffer.String())
}

// ComputeInitialProb ...
func (v *ForwardLog) ComputeInitialProb(pI, pE float64) float64 {
	return pI * pE
}

// ComputeProb ...
func (v *ForwardLog) ComputeProb(tag Tag, index int, pE float64) float64 {
	pSum := 0.0
	for _, givenTag := range v.tags {
		prevResult := (*v.trellis)[givenTag][index - 1]
		prevP := prevResult.Probability
		pT := v.transition.P(tag, givenTag)
		pSum += pT * prevP
	}
	return pE * pSum
}

// Result ...
func (v *ForwardLog) Result(tag Tag, index int) *Result {
	if !v.filled {
		v.FillTrellis()
	}
	return (*v.trellis)[tag][index]
}

// FillTrellis ...
func (v *ForwardLog) FillTrellis() {
	if v.filled {
		return
	}
	for i, _ := range v.sequence {
		v.FillColumn(i)
	}
	v.filled = true
}

// FillColumn ...
func (v *ForwardLog) FillColumn(index int) {
	if v.filled {
		return
	}
	for _, tag := range v.tags {
		(*v.trellis)[tag][index] = v.P(tag, index)
	}
}

// P ...
func (v *ForwardLog) P(tag Tag, index int) *Result {
  value := v.sequence[index]
  pE := v.emission.P(tag, value)
	if index == 0 {
		pI := v.initialState.P(tag)
		p := v.ComputeInitialProb(pI, pE)
		return &Result{"e", p}
	} else {
		p := v.ComputeProb(tag, index, pE)
		return &Result{"e", p}
	}
}
