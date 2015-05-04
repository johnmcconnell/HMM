package hmm

import(
	"fmt"
	"bytes"
	"log"
	"github.com/johnmcconnell/gologspace"
)

type Gamma struct {
	tags []Tag
	sequence []string
	s gologspace.Space
	forward *Forward
	backward *Backward
}

// NewViterb ...
func NewGamma(tags []Tag, sequence []string, s gologspace.Space, i *Initial, t *Transition, e *Emission) *Gamma {
	b := NewBackward(tags, sequence, s, i, t, e)
	b.FillTrellis()
	f := NewForward(tags, sequence, s, i, t, e)
	f.FillTrellis()
	g := Gamma{tags, sequence, s, f, b}
	return &g
}

// String ...
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
		buffer.WriteString(g.R(tag, i).FormatString(g.s))
	}
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}

// ComputeProb ...
func (g *Gamma) CompP(pF, pB, pS float64) float64 {
	s := g.s
	return s.Div(s.Mul(pF, pB), pS)
}

// ComputeProb ...
func (g *Gamma) ComputeP(tag Tag, i int) float64 {
	rF := g.forward.R(tag, i)
	pF := rF.Prob
	rB := g.backward.R(tag, i)
	pB := rB.Prob
	pS := g.CompColumnSum(i)
	return g.CompP(pF, pB, pS)
}

func(g *Gamma) CompTransP(pF, pT, pE, pB float64) float64 {
	s := g.s
	a := s.Mul(s.Mul(pF, pT), s.Mul(pE, pB))
	if (a > 0) {
		log.Printf("a[%v] = pF[%v] + pT[%v] + pE[%v] + pB[%v]", a, pF, pT, pE, pB)
	}
	return a
}

// ComputeTransitionProb ...
func(g *Gamma) ComputeTransP(tag1, tag2 Tag, i int) float64 {
	t := *g.forward.t
	e := *g.forward.e
	value := g.sequence[i + 1]
	pF := g.forward.R(tag1, i).Prob
	pT := t[tag1][tag2]
	pB := g.backward.R(tag2, i + 1).Prob
	pE := e[tag2][value]
	// p := g.CompColumnSum(i)
	return g.CompTransP(pF, pT, pE, pB)
}

// ComputeColumnSum ...
func (g *Gamma) CompColumnSum(i int) float64 {
	pSum := 0.0
	for _, tag := range g.tags {
		rF := g.forward.R(tag, i)
		rB := g.backward.R(tag, i)
		a := g.s.Mul(rF.Prob, rB.Prob)
		if pSum == 0.0 {
			pSum = a
		} else {
			pSum = g.s.Add(a, pSum)
		}
	}
	return pSum
}

// InitialMass ...
func(g *Gamma) InitialMass(tag Tag) float64 {
	return g.ComputeP(tag, 0)
}

// TransitionMass ...
func(g *Gamma) TransMass(tag1, tag2 Tag) float64 {
	pSum := 0.0
	limit := len(g.sequence) - 1
	for i, _ := range g.sequence {
		if (i < limit) {
		  p := g.ComputeTransP(tag1, tag2, i)
			if pSum == 0.0 {
				pSum = p
			} else {
				pSum = g.s.Add(pSum, p)
			}
		}
	}
	// P(tag1, tag2 | w) / P(w)
	return pSum
}

// SumColumn
func (g *Gamma) SumColumn(index int) *Result {
	p := g.CompColumnSum(index)
	return &Result{"e", p}
}

// Result ...
func (g *Gamma) R(tag Tag, index int) *Result {
	p := g.ComputeP(tag, index)
	return &Result{"e", p}
}
