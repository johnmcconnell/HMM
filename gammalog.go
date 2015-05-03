package hmm

import(
	"fmt"
	"bytes"
	"log"
	"github.com/johnmcconnell/gologspace"
)

type GammaLog struct {
	tags []Tag
	sequence []string
	forward *ForwardLog
	backward *BackwardLog
}

// NewViterb ...
func NewGammaLog(tags []Tag, sequence []string, i *InitialState, t *Transition, e *Emission) *GammaLog {
	b := NewBackwardLog(tags, sequence, i, t, e)
	b.FillTrellis()
	f := NewForwardLog(tags, sequence, i, t, e)
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
	pF := rF.Probability
	rB := (*g.backward.trellis)[tag][index]
	pB := rB.Probability
	rSum := g.ComputeColumnSum(index)
	if ((pF + pB) - rSum) > 0 {
		log.Printf("pF[%v] + pB[%v] - rSum[%v]", pF, pB, rSum)
	}
	return (pF + pB) - rSum
}

// ComputeTransitionProb ...
func(g *GammaLog) ComputeTransitionProb(tag1, tag2 Tag, i int) float64 {
	t := *g.forward.transition
	e := *g.forward.emission
	value := g.sequence[i + 1]
	pF := (*g.forward.trellis)[tag1][i].Probability
	pT := gologspace.LogProb(t[tag1][tag2])
	pB := (*g.backward.trellis)[tag2][i + 1].Probability
	pE := gologspace.LogProb(e[tag2][value])
	p := g.ComputeColumnSum(i)
	a := pF + pT + pB + pE // - p
	if (a > 0) {
		log.Printf("a[%v] = pF[%v] + pT[%v] + pB[%v] + pE[%v] - p[%v]", a, pF, pT, pB, pE, p)
	}
	return a
}

// ComputeColumnSum ...
func (g *GammaLog) ComputeColumnSum(index int) float64 {
	pSum := 0.0
	for _, tag := range g.tags {
		rF := (*g.forward.trellis)[tag][index]
		rB := (*g.backward.trellis)[tag][index]
		a := rF.Probability + rB.Probability
		if pSum == 0.0 {
			pSum = a
		} else {
			pSum = gologspace.LogAdd(a, pSum)
		}
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
			if pSum == 0.0 {
				pSum = p
			} else {
				pSum = gologspace.LogAdd(pSum, p)
			}
		}
	}
	// P(tag1, tag2 | w) / P(w)
	return pSum
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
