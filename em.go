package hmm

type EM struct{
	tags []Tag
	sentences [][]string
	i InitialState
	t Transition
	e Emission
}

func NewEM(tags []Tag, sentences [][]string,
i InitialState, t Transition, e Emission) *EM {
	return &EM{tags, sentences, i, t, e}
}

func (e *EM) String() string {
	return "Hello"
}

func (e *EM) Next() *EM {
	eNext := *e
	iP, tP, eP := e.EStep()
	e.MStep(iP, tP, eP)
	return &eNext
}

func (e *EM) EStep() (*InitialState, *Transition, *Emission) {
	iMass := make(InitialState)
	eMass := NewEmission(e.tags)
	tMass := NewTransition(e.tags)
	for _, sentence := range e.sentences {
		g := NewGamma(e.tags, sentence, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {
			iMass[tag] = g.InitialMass(tag)
			for _, tag2 := range e.tags {
				tMass[tag][tag2] = g.TransitionMass(tag, tag2)
			}
			for i, word := range sentence {
				eMass[tag][word] = g.ComputeProb(tag, i)
			}
		}
	}
	return &iMass, &tMass, &eMass
}

func (e *EM) MStep(iP *InitialState, tP *Transition, eP *Emission) { }
