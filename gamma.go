package hmm

import(
	"fmt"
	"bytes"
)

type GammaLog struct {
	tags []Tag
	sequence []string
	forward *Forward
	backward *Backward
}

// NewViterb ...
func NewGammaLog(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *GammaLog {
	b := NewBackward(tags, sequence, i, t, e)
	b.FillTrellis()
	f := NewForward(tags, sequence, i, t, e)
	f.FillTrellis()
	g := GammaLog{tags, sequence, f, b}
	return &g
}

// String ...
func (g *GammaLog) String() string {
	buffer := bytes.NewBufferString(fmt.Sprintf("GammaLog: '%s'\n", g.sequence))
	for _, tag := range g.tags {
		buffer.WriteString(g.RowString(tag))
	}
	buffer.WriteString(g.SumRowString())
	return fmt.Sprintf("%s\n", buffer.String())
}

// ComputeProb ...
func (g *GammaLog) ComputeProb(tag Tag, index int) float64 {
	rF := (*g.forward.trellis)[tag][index]
	rB := (*g.backward.trellis)[tag][index]
	rSum := g.SumColumn(index)
	return (rF.Probability * rB.Probability) / rSum.Probability
}

// ComputeTransitionProb ...
func(g *GammaLog) ComputeTransitionProb(tag1, tag2 Tag, i int) float64 {
	t := *g.forward.transition
	e := *g.forward.emission
	value := g.sequence[i]
	pF := (*g.forward.trellis)[tag1][i].Probability
	pT := t[tag2][tag1]
	pB := (*g.backward.trellis)[tag2][i + 1].Probability
	pE := e[tag2][value]
	return pF * pT * pB * pE
}

// ComputeColumnSum ...
func (g *GammaLog) ComputeColumnSum(index int) float64 {
	pSum := 0.0
	for _, tag := range g.tags {
		rF := (*g.forward.trellis)[tag][index]
		rB := (*g.backward.trellis)[tag][index]
		pSum += (rF.Probability * rB.Probability)
	}
	return pSum
}

// InitialMass ...
func(g *GammaLog) InitialMass(tag Tag) float64 {
	return g.ComputeProb(tag, 0)
}

// TransitionMass ...
func(g *GammaLog) TransitionMass(tag1, tag2 Tag) float64 {
	pSum := 0.0
	limit := len(g.sequence) - 1
	for i, _ := range g.sequence {
		if (i < limit) {
		  p := g.ComputeTransitionProb(tag1, tag2, i)
		  pSum += p
		}
	}
	// P(tag1, tag2 | w) / P(w)
	return pSum / g.ComputeSentenceProb()
}

// ComputeSentenceProb ...
func(g *GammaLog) ComputeSentenceProb() float64 {
	return g.ComputeColumnSum(0)
}

// SumRowString ...
func (g *GammaLog) SumRowString() string {
	buffer := bytes.NewBufferString("|Sum: ' '|")
	for i, _ := range g.sequence {
		r := g.SumColumn(i)
		buffer.WriteString(r.String())
	}
	return buffer.String()
}

// RowString ...
func (g *GammaLog) RowString(tag Tag) string {
	buffer := bytes.NewBufferString(fmt.Sprintf("|Tag: '%v'|", tag))
	for i, _ := range g.sequence {
		buffer.WriteString(g.Result(tag, i).String())
	}
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}

// SumColumn
func (g *GammaLog) SumColumn(index int) *Result {
	p := g.ComputeColumnSum(index)
	return &Result{"e", p}
}

// Result ...
func (g *GammaLog) Result(tag Tag, index int) *Result {
	p := g.ComputeProb(tag, index)
	return &Result{tag, p}
}
