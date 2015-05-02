package hmm

import(
	"fmt"
	"bytes"
)

type Gamma struct {
	tags []Tag
	sequence string
	forward *Forward
	backward *Backward
}

// NewViterb ...
func NewGamma(tags []Tag, sequence string, i *InitialState, t *Transition, e *Emission) *Gamma {
	b := NewBackward(tags, sequence, i, t, e)
	b.FillTrellis()
	f := NewForward(tags, sequence, i, t, e)
	f.FillTrellis()
	g := Gamma{tags, sequence, f, b}
	return &g
}

func (g *Gamma) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("Gamma: '%s'\n", g.sequence))
	for _, tag := range g.tags {
		buffer.WriteString(g.RowString(tag))
	}
	buffer.WriteString(g.SumRowString())
	return fmt.Sprintf(buffer.String())
}

// SumRowString ...
func (g *Gamma) SumRowString() string {
	buffer := bytes.NewBufferString("|Sum: ' '|")
	for i, _ := range g.sequence {
		r := &Result{"e", 0.0}
		for _, tag := range g.tags {
		  r.Probability += g.Result(tag, i).Probability
		}
		buffer.WriteString(r.String())
	}
	return buffer.String()
}

// RowString ...
func (g *Gamma) RowString(tag Tag) string {
	buffer := bytes.NewBufferString(fmt.Sprintf("|Tag: '%v'|", tag))
	for i, _ := range g.sequence {
		buffer.WriteString(g.Result(tag, i).String())
	}
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}

// Result ...
func (v *Gamma) Result(tag Tag, index int) *Result {
	rF := (*v.forward.trellis)[tag][index]
	rB := (*v.backward.trellis)[tag][index]
	return &Result{tag, rF.Probability * rB.Probability}
}
