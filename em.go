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
	// e.Check(iP, tP, eP)
	//e.Apply(iP, tP, eP, e.s.Enter)
	return NewEM(e.tags, e.words, e.sentences, e.s, *iP, *tP, *eP)
}

func (e *EM) EStep() (*Initial, *Transition, *Emission, *Initial, *Initial) {
	iCount := make(Initial)
	eCount := NewEmission(e.tags)
	tCount := NewTransition(e.tags)
	tagCount := make(Initial)
	tagTagCount := make(Initial)

	for iS, sentence := range e.sentences {
		g := NewGamma(e.tags, sentence, e.s, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {

			iTMass := g.InitialMass(tag)
			// Must check if unitialized or log add won't work.
			if iCount[tag] == 0.0 {
				iCount[tag] = iTMass
			} else {
				iCount[tag] = e.s.Add(iCount[tag], iTMass)
			}

			for _, tag2 := range e.tags {
				tTMass := g.TransMass(tag, tag2)
				// Must check if unitialized or log add won't work.
				if tCount[tag][tag2] == 0.0 {
					tCount[tag][tag2] = tTMass
				} else {
					tCount[tag][tag2] = e.s.Add(tCount[tag][tag2], tTMass)
				}
			}

			lastWordIndex := len(sentence) - 1
			for iW, word := range sentence {
				eMass := g.ComputeP(tag, iW)
				// Must check if unitialized or log add won't work.
				if eCount[tag][word] == 0.0 {
					eCount[tag][word] = eMass
				} else {
					eCount[tag][word] = e.s.Add(eCount[tag][word], eMass)
				}
				// Must check if unitialized or log add won't work.
				if tagCount[tag] == 0.0 {
					tagCount[tag] = eMass
				} else {
					tagCount[tag] = e.s.Add(tagCount[tag], eMass)
				}

				// All but last word in sentence
				if (iW < lastWordIndex) {
					if tagTagCount[tag] == 0.0 {
						tagTagCount[tag] = eMass
					} else {
						tagTagCount[tag] = e.s.Add(tagTagCount[tag], eMass)
					}
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
	sLength := e.s.Enter(float64(len(e.sentences)))

	for _, tag := range e.tags {
		(*iP)[tag] = e.s.Div((*iP)[tag], sLength)

		tagCount := (*tC)[tag]
		tagToTagCount := (*ttC)[tag]

		for _, tag2 := range e.tags {
			(*tP)[tag][tag2] = e.s.Div((*tP)[tag][tag2], tagToTagCount)
		}

		for word, _ := range (*eP)[tag] {
			(*eP)[tag][word] = e.s.Div((*eP)[tag][word], tagCount)
		}
	}
}

func (e *EM) Apply(iP *Initial, tP *Transition, eP *Emission, f func (float64) float64) {
	for _, tag := range e.tags {
		(*iP)[tag] = f((*iP)[tag])
		for _, tag2 := range e.tags {
			(*tP)[tag][tag2] = f((*tP)[tag][tag2])
		}
		for _, word := range e.words {
			(*eP)[tag][word] = f((*eP)[tag][word])
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
