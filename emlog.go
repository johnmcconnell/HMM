package hmm

import(
	"math"
	"log"
	"os"
)

type EMLog struct{
	tags []Tag
	sentences [][]string
	i InitialState
	t Transition
	e Emission
}

func NewEMLog(tags []Tag, sentences [][]string,
i InitialState, t Transition, e Emission) *EMLog {
	return &EMLog{tags, sentences, i, t, e}
}

func NewEMLog2(tags []Tag, sentences [][]string,
i InitialState, t Transition, e Emission) *EMLog {
	em := EMLog{tags, sentences, i, t, e}
	em.Check(&i, &t, &e)
	return &em
}

func (e *EMLog) Next() *EMLog {
	iP, tP, eP, tC := e.EStep()
	e.MStep(iP, tC, tP, eP)
	return NewEMLog(e.tags, e.sentences, *iP, *tP, *eP)
}

func (e *EMLog) EStep() (*InitialState, *Transition, *Emission, *InitialState) {
	iCount := make(InitialState)
	tagCount := make(InitialState)
	eCount := NewEmission(e.tags)
	tCount := NewTransition(e.tags)
	for _, sentence := range e.sentences {
		g := NewGammaLog(e.tags, sentence, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {
			iCount[tag] += math.Exp(g.InitialMass(tag))
			for _, tag2 := range e.tags {
				tCount[tag][tag2] += math.Exp(g.TransitionMass(tag, tag2))
			}
			for i, word := range sentence {
				p := math.Exp(g.ComputeProb(tag, i))
				eCount[tag][word] += p
				tagCount[tag] += p
			}
		}
	}
	return &iCount, &tCount, &eCount, &tagCount
}

func (e *EMLog) Check(iP *InitialState, tP *Transition, eP *Emission) {
	for _, sentence := range e.sentences {
		for _, tag := range e.tags {
				iP := (*iP)[tag]
				if (iP > 1.0) {
					log.Printf("Invalid 'I(%s) = '%v'", tag, iP)
					os.Exit(-1)
				}
			for _, tag2 := range e.tags {
				tP := (*tP)[tag][tag2]
				if (tP > 1.0) {
					log.Printf("Invalid 'T(%s|%s) = '%v'", tag2, tag, tP)
					os.Exit(-1)
				}
			}
			for _, word := range sentence {
				eP := (*eP)[tag][word]
				if (eP > 1.0) {
					log.Printf("Invalid 'E(%s|%s) = '%v'", word, tag, )
					os.Exit(-1)
				}
			}
		}
	}
}

func (e *EMLog) MStep(iP, tC *InitialState, tP *Transition, eP *Emission) {
	lS := float64(len(e.sentences))
	for _, tag := range e.tags {
		p := (*iP)[tag] / lS
		if (p > 1.0) {
			log.Printf("IP(%s) = %v = (%v / %v)", tag, p, (*iP)[tag], lS)
		}
		(*iP)[tag] = p
		tagCount := (*tC)[tag]
		for _, tag2 := range e.tags {
			p = (*tP)[tag][tag2] / tagCount
			if (p > 1.0) {
				log.Printf("TP(%s|%s) = %v = (%v / %v)", tag2, tag, p, (*tP)[tag][tag2], tagCount)
			}
			(*tP)[tag][tag2] = p
		}
		for word, _ := range (*eP)[tag] {
			p = (*eP)[tag][word] / tagCount
			if (p > 1.0) {
				log.Printf("EP(%s|%s) = %v = (%v / %v)", word, tag, p, (*eP)[tag][word], tagCount)
			}
			(*eP)[tag][word] = p
		}
	}
	e.Check(iP, tP, eP)
}

func (e *EMLog) I() InitialState { return e.i }
func (e *EMLog) T() Transition { return e.t }
func (e *EMLog) E() Emission { return e.e }
