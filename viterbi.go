package hmm

import(
	"fmt"
	"bytes"
	"strings"
)

type Viterbi struct {
	tags []Tag
	sequence string
	filled bool
	trellis *Trellis
	initialState *InitialState
	transition *Transition
	emission *Emission
}

// NewViterbi ...
func NewViterbi(tags []Tag, sequence string, i *InitialState, t *Transition, e *Emission) *Viterbi {
	trellis := NewTrellis(tags, len(sequence))
	v := Viterbi{tags, sequence, false, trellis, i, t, e}
	return &v
}

func (v *Viterbi) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Viterbi: '%s'\n", v.sequence))
	buffer.WriteString(v.trellis.String())
	p := v.Prediction()
	buffer.WriteString(fmt.Sprintf("Prediction: %s", p))
	return fmt.Sprintf(buffer.String())
}

// Result ...
func (v *Viterbi) Result(index int) *Result {
	if !v.filled {
		v.FillTrellis()
	}
	return v.MaxResult(index)
}

// FillTrellis ...
func (v *Viterbi) FillTrellis() {
	if v.filled {
		return
	}
	for i, _ := range v.sequence {
		v.FillColumn(i)
	}
	v.filled = true
}

// FillColumn ...
func (v *Viterbi) FillColumn(index int) {
	if v.filled {
		return
	}
	for _, tag := range v.tags {
		(*v.trellis)[tag][index] = v.P(tag, index)
	}
}

// P ...
func (v *Viterbi) P(tag Tag, index int) *Result {
	if index == 0 {
		value := v.sequence[index]
		pI := v.initialState.P(tag)
		pE := v.emission.P(tag, value)
		return &Result{"e", pI * pE}
	} else {
	  return v.MaxP(tag, index)
	}
}

// MaxP ...
func (v *Viterbi) MaxP(tag Tag, index int) *Result {
	var maxResult *Result = nil
  value := v.sequence[index]
	for _, givenTag := range v.tags {
		prevResult := (*v.trellis)[givenTag][index - 1]
		pT := v.transition.P(tag, givenTag)
		pE := v.emission.P(tag, value)
		p := prevResult.Probability * pT * pE
		if (maxResult == nil) {
			maxResult = &Result{givenTag, p}
		}
		if (maxResult.Probability < p) {
			maxResult = &Result{givenTag, p}
		}
	}
	return maxResult
}

// MaxResult ...
func (v *Viterbi) MaxResult(index int) *Result {
	var maxResult *Result = nil
	for _, tag := range v.tags {
		currResult := (*v.trellis)[tag][index]
		if (maxResult == nil) {
			maxResult = currResult
		}
		if (maxResult.Probability < currResult.Probability) {
			maxResult = currResult
		}
	}
	return maxResult
}

// MaxTag ...
func (v *Viterbi) MaxTag(index int) Tag {
	var maxTag Tag = ""
	var maxProb float64 = 0.0
	for _, tag := range v.tags {
		currResult := (*v.trellis)[tag][index]
		if (maxProb < currResult.Probability) {
			maxTag = tag
			maxProb = currResult.Probability
		}
	}
	return maxTag
}


// Prediction ...
func (v *Viterbi) Prediction() string {
	l := len(v.sequence)
	previousTag := v.MaxTag(l - 1)
	tags := []string{string(previousTag)}
	for i, _ := range v.sequence {
		reverseI := l - i - 1
		r := (*v.trellis)[previousTag][reverseI]
		previousTag = r.previousTag
		tags = append([]string{string(previousTag)}, tags...)
	}
	return strings.Join(tags, "")
}
