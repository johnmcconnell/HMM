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
	tagMass := make(InitialState)
	transitionMass := NewTransition(e.tags)
	for _, sentence := range e.sentences {
		g := NewGamma(e.tags, sentence, &e.i, &e.t, &e.e)
		for _, tag := range e.tags {
			tagMass[tag] = g.TagMass(tag)
			for _, tag2 := range e.tags {
				transitionMass[tag][tag2] = g.TransitionMass(tag, tag2)
			}
		}
	}
	return &eNext
}
