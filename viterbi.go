package hmm

import(
	"fmt"
	"bytes"
	"strings"
	"github.com/johnmcconnell/gologspace"
)

type Viterbi struct {
	tags []Tag
	sequence []string
	cache *Trellis
	s gologspace.Space
	i *Initial
	t *Transition
	e *Emission
}

// NewViterbi ...
func NewViterbi(tags []Tag, sequence []string, s gologspace.Space, i *Initial, t *Transition, e *Emission) *Viterbi {
	cache := NewTrellis(tags, len(sequence))
	v := Viterbi{tags, sequence, cache, s, i, t, e}
	return &v
}

// String ...
func (v *Viterbi) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Viterbi: '%s'\n", v.sequence))
	buffer.WriteString(v.cache.FormatString(v.s))
	p := v.Prediction()
	buffer.WriteString(fmt.Sprintf("Prediction: %s\n", p))
	l, _ := v.Labeled()
	buffer.WriteString(fmt.Sprintf("Label: %s\n", l))
	return fmt.Sprintf(buffer.String())
}

// ComputeInitialProb ...
func (v *Viterbi) CompInitialP(pI, pE float64) float64 {
	return v.s.Mul(pI, pE)
}

// ComputeProb ...
func (v *Viterbi) CompP(prevMax, pE, pT float64) float64 {
	return v.s.Mul(v.s.Mul(prevMax, pT), pE)
}

// FillTrellis ...
func (v *Viterbi) FillTrellis() {
	for i, _ := range v.sequence {
		v.FillColumn(i)
	}
}

// FillColumn ...
func (v *Viterbi) FillColumn(index int) {
	for _, tag := range v.tags {
		v.FillValue(tag, index)
	}
}

// P ...
func (v *Viterbi) FillValue(tag Tag, i int) {
	val := v.sequence[i]
	pE := v.e.P(tag, val)
	if i == 0 {
		pI := v.i.P(tag)
		p := v.CompInitialP(pI, pE)
		(*v.cache)[tag][i] = &Result{"e", p}
	} else {
		givenTag := v.cache.MaxColumn(i - 1)
		prevR := v.R(givenTag, i - 1)
		prevP := prevR.Prob
		pT := v.t.P(tag, givenTag)
		p := v.CompP(prevP, pE, pT)
		(*v.cache)[tag][i] = &Result{givenTag, p}
	}
}

// Prediction ...
func (v *Viterbi) Prediction() string {
	l := len(v.sequence)
	t := v.cache.MaxColumn(l - 1)
	prevTag := v.R(t, l - 1).prevTag
	tags := []string{string(prevTag)}
	for i, _ := range v.sequence {
		rI := l - i - 1
		r := (*v.cache)[prevTag][rI]
		prevTag = r.prevTag
		tags = append([]string{string(prevTag)}, tags...)
	}
	return strings.Join(tags, "")
}

// Labeled ...
func (v *Viterbi) Labeled() ([]LabeledWord, error) {
	l := len(v.sequence) - 1
	tag := v.cache.MaxColumn(l)
	word := v.sequence[l]
	label := LabeledWord{word, tag}
	tags := []LabeledWord{label}
	for i, _ := range v.sequence {
		if (i < l) {
			rI := l - i
			word = v.sequence[rI - 1]
			q := (*v.cache)[tag]
			r := q[rI]
			tag = r.prevTag
			label := LabeledWord{word, tag}
			tags = append([]LabeledWord{label}, tags...)
		}
	}
	return tags, nil
}

// R ...
func (v *Viterbi) R(t Tag, i int) *Result {
	return (*v.cache)[t][i]
}
