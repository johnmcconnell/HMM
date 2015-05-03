package hmm

import(
	"math"
	"log"
	"os"
)

type EMLog struct{
	tags []Tag
	words []string
	sentences [][]string
	i InitialState
	t Transition
	e Emission
}

func NewEMLog(tags []Tag, words []string, sentences [][]string,
i InitialState, t Transition, e Emission) *EMLog {
	return &EMLog{tags, words, sentences, i, t, e}
}

func NewEMLog2(tags []Tag, words []string, sentences [][]string,
i InitialState, t Transition, e Emission) *EMLog {
	em := EMLog{tags, words, sentences, i, t, e}
	em.Check(&i, &t, &e)
	return &em
}

func (e *EMLog) Next() *EMLog {
	iP, tP, eP, tC, ttC := e.EStep()
	e.MStep(iP, tC, ttC, tP, eP)
	return NewEMLog(e.tags, e.words, e.sentences, *iP, *tP, *eP)
}

func (e *EMLog) EStep() (*InitialState, *Transition, *Emission, *InitialState, *InitialState) {
	iCount := make(InitialState)
	tagCount := make(InitialState)
	tagTagCount := make(InitialState)
	eCount := NewEmission(e.tags)
	tCount := NewTransition(e.tags)
	for iS, sentence := range e.sentences {
		g := NewGammaLog(e.tags, sentence, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {
			iCount[tag] += math.Exp(g.InitialMass(tag))
			for _, tag2 := range e.tags {
				count := math.Exp(g.TransitionMass(tag, tag2))
				tCount[tag][tag2] += count
			}
			limit := len(sentence) - 1
			for iW, word := range sentence {
				p := math.Exp(g.ComputeProb(tag, iW))
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
	iPT := 0.0
	for _, tag := range e.tags {
		tPT := 0.0
		for _, tag2 := range e.tags {
			tPT += (*tP)[tag][tag2]
		}
		d := math.Abs(1.0 - tPT)
		if (d > 0.01) {
			log.Printf("Invalid 'TP(%s) Total = '%v'", tag, tPT)
		}
		ePT := 0.0
		for _, word := range e.words {
			ePT += (*eP)[tag][word]
		}
		d = math.Abs(1.0 - ePT)
		if (d > 0.01) {
			log.Printf("Invalid 'EP(%s) Total = '%v'", tag, ePT)
		}
	}
	if (iPT > 1.0) {
		log.Printf("Invalid 'IP Total = '%v'", iPT)
	}
}

func (e *EMLog) MStep(iP, tC, ttC *InitialState, tP *Transition, eP *Emission) {
	lS := float64(len(e.sentences))
	for _, tag := range e.tags {
		(*iP)[tag] = (*iP)[tag] / lS
		tagCount := (*tC)[tag]
		tagTagCount := (*ttC)[tag]
		for _, tag2 := range e.tags {
			(*tP)[tag][tag2] = (*tP)[tag][tag2] * 67.7822993789443 / tagTagCount
		}
		for word, _ := range (*eP)[tag] {
			(*eP)[tag][word] = (*eP)[tag][word] / tagCount
		}
	}
	e.Check(iP, tP, eP)
}

func (e *EMLog) I() InitialState { return e.i }
func (e *EMLog) T() Transition { return e.t }
func (e *EMLog) E() Emission { return e.e }
