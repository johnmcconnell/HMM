package hmm

type EM struct{
	sentences [][]string
	i InitialState
	t Transition
	e Emission
}

func NewEM(sentences [][]string,
i InitialState, t Transition, e Emission) *EM {
	return &EM{sentences, i, t, e}
}

func (e *EM) String() string {
	return "Hello"
}
