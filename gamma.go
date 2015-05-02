package hmm

import(
	"fmt"
	"bytes"
)

type Gamma struct {
	tags []Tag
	sequence []string
	forward *Forward
	backward *Backward
}

// NewViterb ...
func NewGamma(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *Gamma {
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
	return fmt.Sprintf("%s\n", buffer.String())
}

// SumRowString ...
func (g *Gamma) SumRowString() string {
	buffer := bytes.NewBufferString("|Sum: ' '|")
	for i, _ := range g.sequence {
		r := g.SumColumn(i)
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

// SumColumn
func (g *Gamma) SumColumn(index int) *Result {
	prob := 0.0
	for _, tag := range g.tags {
		rF := (*g.forward.trellis)[tag][index]
		rB := (*g.backward.trellis)[tag][index]
		prob += (rF.Probability * rB.Probability)
	}
	return &Result{"e", prob}
}

// Result ...
func (g *Gamma) Result(tag Tag, index int) *Result {
	rF := (*g.forward.trellis)[tag][index]
	rB := (*g.backward.trellis)[tag][index]
	rSum := g.SumColumn(index)
	prob := (rF.Probability * rB.Probability) / rSum.Probability
	return &Result{tag, prob}
}
