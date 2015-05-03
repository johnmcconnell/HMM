package hmm

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

func (e *EMLog) String() string {
	return "Hello"
}

func (e *EMLog) Next() *EMLog {
	eNext := *e
	iP, tP, eP, tC := e.EStep()
	e.MStep(iP, tC, tP, eP)
	return &eNext
}

func (e *EMLog) EStep() (*InitialState, *Transition, *Emission, *InitialState) {
	iCount := make(InitialState)
	tagCount := make(InitialState)
	eCount := NewEmission(e.tags)
	tCount := NewTransition(e.tags)
	for _, sentence := range e.sentences {
		g := NewGamma(e.tags, sentence, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {
			iCount[tag] += g.InitialMass(tag)
			for _, tag2 := range e.tags {
				tCount[tag][tag2] += g.TransitionMass(tag, tag2)
			}
			for i, word := range sentence {
				p := g.ComputeProb(tag, i)
				eCount[tag][word] += p
				tagCount[tag] += p
			}
		}
	}
	return &iCount, &tCount, &eCount, &tagCount
}

func (e *EMLog) MStep(iP, tC *InitialState, tP *Transition, eP *Emission) {
	lS := len(e.sentences)
	for _, tag := range e.tags {
		(*iP)[tag] = (*iP)[tag] / float64(lS)
		tagCount := (*tC)[tag]
		for _, tag2 := range e.tags {
			(*tP)[tag][tag2] = (*tP)[tag][tag2] / tagCount
		}
		for word, _ := range (*eP)[tag] {
			(*eP)[tag][word] = (*eP)[tag][word] / tagCount
		}
	}
}

func (e *EMLog) I() InitialState { return e.i }
func (e *EMLog) T() Transition { return e.t }
func (e *EMLog) E() Emission { return e.e }
