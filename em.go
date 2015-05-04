package hmm

import(
	"math"
	"log"
	"os"
	"github.com/johnmcconnell/gologspace"
)

type EM struct{
	tags []Tag
	words []string
	sentences [][]string
	s gologspace.Space
	i Initial
	t Transition
	e Emission
}

func NewEM(tags []Tag, words []string, sentences [][]string,
s gologspace.Space, i Initial, t Transition, e Emission) *EM {
	return &EM{tags, words, sentences, s, i, t, e}
}

func NewEM2(tags []Tag, words []string, sentences [][]string,
s gologspace.Space, i Initial, t Transition, e Emission) *EM {
	em := NewEM(tags, words, sentences, s, i, t, e)
	return em
}

func (e *EM) Next() *EM {
	iP, tP, eP, tC, ttC := e.EStep()
	e.MStep(iP, tC, ttC, tP, eP)
	e.Check(iP, tP, eP)
	e.ReenterSpace(iP, tP, eP)
	return NewEM(e.tags, e.words, e.sentences, e.s, *iP, *tP, *eP)
}

func (e *EM) EStep() (*Initial, *Transition, *Emission, *Initial, *Initial) {
	iCount := make(Initial)
	tagCount := make(Initial)
	tagTagCount := make(Initial)
	eCount := NewEmission(e.tags)
	tCount := NewTransition(e.tags)
	for iS, sentence := range e.sentences {
		g := NewGamma(e.tags, sentence, e.s, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {
			iCount[tag] += e.s.Exit(g.InitialMass(tag))
			for _, tag2 := range e.tags {
				count := e.s.Exit(g.TransMass(tag, tag2))
				tCount[tag][tag2] += count
			}
			limit := len(sentence) - 1
			for iW, word := range sentence {
				p := e.s.Exit(g.ComputeP(tag, iW))
				eCount[tag][word] += p
			  tagCount[tag] += p
				if (iW < limit) {
					tagTagCount[tag] += p
				}
			}
		}
		if (iS % 3000 == 0) {
			log.Printf("3000 Finished in E Step\n")
		}
	}
	return &iCount, &tCount, &eCount, &tagCount, &tagTagCount
}

func (e *EM) MStep(iP, tC, ttC *Initial, tP *Transition, eP *Emission) {
	lS := float64(len(e.sentences))
	for _, tag := range e.tags {
		(*iP)[tag] = (*iP)[tag] / lS
		tagCount := (*tC)[tag]
		tagTagCount := (*ttC)[tag]
		for _, tag2 := range e.tags {
			(*tP)[tag][tag2] = (*tP)[tag][tag2] /  tagTagCount
		}
		for word, _ := range (*eP)[tag] {
			(*eP)[tag][word] = (*eP)[tag][word] / tagCount
		}
	}
}

func (e *EM) ReenterSpace(iP *Initial, tP *Transition, eP *Emission) {
	for _, tag := range e.tags {
		(*iP)[tag] = e.s.Enter((*iP)[tag])
		for _, tag2 := range e.tags {
			(*tP)[tag][tag2] = e.s.Enter((*tP)[tag][tag2])
		}
		for _, word := range e.words {
			(*eP)[tag][word] = e.s.Enter((*eP)[tag][word])
		}
	}
}

func (e *EM) Check(iP *Initial, tP *Transition, eP *Emission) {
	one := e.s.Enter(1.0)
	for _, sentence := range e.sentences {
		for _, tag := range e.tags {
				iP := (*iP)[tag]
				if (iP > one) {
					log.Printf("Invalid 'I(%s) = '%v'", tag, iP)
					os.Exit(-1)
				}
			for _, tag2 := range e.tags {
				tP := (*tP)[tag][tag2]
				if (tP > one) {
					log.Printf("Invalid 'T(%s|%s) = '%v'", tag2, tag, tP)
					os.Exit(-1)
				}
			}
			for _, word := range sentence {
				eP := (*eP)[tag][word]
				if (eP > one) {
					log.Printf("Invalid 'E(%s|%s) = '%v'", word, tag, )
					os.Exit(-1)
				}
			}
		}
	}
	iPT := 0.0
	for _, tag := range e.tags {
		tPT := 0.0
		for _, tag2 := range e.tags {
			tPT += (*tP)[tag][tag2]
		}
		d := math.Abs(e.s.Sub(one, tPT))
		if (d > 0.01) {
			log.Printf("Invalid 'TP(%s) Total = '%v'", tag, tPT)
		}
		ePT := 0.0
		for _, word := range e.words {
			ePT += (*eP)[tag][word]
		}
		d = math.Abs(e.s.Sub(one, ePT))
		if (d > 0.01) {
			log.Printf("Invalid 'EP(%s) Total = '%v'", tag, ePT)
		}
	}
	if (iPT > one) {
		log.Printf("Invalid 'IP Total = '%v'", iPT)
	}
}

func (e *EM) I() Initial { return e.i }
func (e *EM) T() Transition { return e.t }
func (e *EM) E() Emission { return e.e }
