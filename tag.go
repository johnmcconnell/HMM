package hmm

type Tag string

func (t *Tag) Blank() bool {
	return (*t == "")
}
