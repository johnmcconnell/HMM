package hmm

type EM struct{}

func NewEM(sentences [][]string) *EM {
	return &EM{}
}

func (e *EM) String() string {
	return "Hello"
}
