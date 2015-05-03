package hmm

import(
	"fmt"
	"bytes"
	"errors"
	"strings"
	"math"
	"log"
	"github.com/johnmcconnell/gologspace"
)

type ViterbiLog struct {
	tags []Tag
	sequence []string
	filled bool
	trellis *Trellis
	initialState *InitialState
	transition *Transition
	emission *Emission
}

// NewViterbiLog ...
func NewViterbiLog(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *ViterbiLog {
	trellis := NewTrellis(tags, len(sequence))
	v := ViterbiLog{tags, sequence, false, trellis, i, t, e}
	return &v
}

// String ...
func (v *ViterbiLog) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("ViterbiLog: '%s'\n", v.sequence))
	buffer.WriteString(v.trellis.String())
	p := v.Prediction()
	buffer.WriteString(fmt.Sprintf("Prediction: %s\n", p))
	l, _ := v.Labeled()
	buffer.WriteString(fmt.Sprintf("Label: %s\n", l))
	return fmt.Sprintf(buffer.String())
}

// ComputeInitialProb ...
func (v *ViterbiLog) ComputeInitialProb(pI, pE float64) float64 {
	return pI + pE
}

// ComputeProb ...
func (v *ViterbiLog) ComputeProb(prevMax, pE, pT float64) float64 {
	return prevMax + pT + pE
}

// Result ...
func (v *ViterbiLog) Result(index int) *Result {
	if !v.filled {
		v.FillTrellis()
	}
	return v.MaxResult(index)
}

// FillTrellis ...
func (v *ViterbiLog) FillTrellis() {
	if v.filled {
		return
	}
	for i, _ := range v.sequence {
		v.FillColumn(i)
	}
	v.filled = true
}

// FillColumn ...
func (v *ViterbiLog) FillColumn(index int) {
	if v.filled {
		return
	}
	for _, tag := range v.tags {
		(*v.trellis)[tag][index] = v.P(tag, index)
	}
}

// P ...
func (v *ViterbiLog) P(tag Tag, index int) *Result {
	if index == 0 {
		value := v.sequence[index]
		pI := gologspace.LogProb(v.initialState.P(tag))
		pE := gologspace.LogProb(v.emission.P(tag, value))
		p := v.ComputeInitialProb(pI, pE)
		return &Result{"e", p}
	} else {
	  return v.MaxP(tag, index)
	}
}

// MaxP ...
func (v *ViterbiLog) MaxP(tag Tag, index int) *Result {
	var maxResult *Result = nil
  value := v.sequence[index]
	for _, givenTag := range v.tags {
		prevResult := (*v.trellis)[givenTag][index - 1]
		pT := gologspace.LogProb(v.transition.P(tag, givenTag))
		pE := gologspace.LogProb(v.emission.P(tag, value))
		p := v.ComputeProb(prevResult.Probability, pT, pE)
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
func (v *ViterbiLog) MaxResult(index int) *Result {
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
func (v *ViterbiLog) MaxTag(index int) Tag {
	var maxTag Tag = ""
	maxProb := math.Inf(-1)
	for _, tag := range v.tags {
		currResult := (*v.trellis)[tag][index]
		if (maxProb <= currResult.Probability) {
			maxTag = tag
			maxProb = currResult.Probability
		}
	}
	return maxTag
}

// Prediction ...
func (v *ViterbiLog) Prediction() string {
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

func (v *ViterbiLog) Labeled() ([]LabeledWord, error) {
	l := len(v.sequence) - 1
	tag := v.MaxTag(l)
	if (tag == "") {
		log.Println("Failed on First")
	}
	word := v.sequence[l]
	label := LabeledWord{word, tag}
	tags := []LabeledWord{label}
	for i, _ := range v.sequence {
		if (i < l) {
			reverseI := l - i
			word = v.sequence[reverseI - 1]
			q := (*v.trellis)[tag]
			if (reverseI >= len(q)) {
				s := fmt.Sprintf("On '%s': i:%v vs len:%v, q[U] = %v\n", tag,
				reverseI, len(q), (*v.trellis)["UNKNOWN"])
				return nil, errors.New(s)
			} else {
				r := q[reverseI]
				tag = r.previousTag
				label := LabeledWord{word, tag}
				tags = append([]LabeledWord{label}, tags...)
			}
		}
	}
	return tags, nil
}
