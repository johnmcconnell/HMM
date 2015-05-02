package hmm

import(
	"fmt"
	"bytes"
)

type Result struct {
	previousTag Tag
	Probability float64
}

type Trellis map[Tag][]*Result

// String ...
func (t *Trellis) String() string {
	buffer := bytes.NewBufferString("Trellis:\n")
	for tag, results := range *t {
		buffer.WriteString(t.RowString(tag, results))
	}
	return fmt.Sprintf(buffer.String())
}

// RowString ...
func (t *Trellis) RowString(tag Tag, results []*Result) string {
	buffer := bytes.NewBufferString(fmt.Sprintf("|Tag: '%v'|", tag))
	for _, result := range results {
		buffer.WriteString(result.String())
	}
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}

// String ...
func (r *Result) String() string {
	return fmt.Sprintf("'%v': %.6fp|", r.previousTag, r.Probability)
}

// New ...
func NewTrellis(tags []Tag, size int) *Trellis {
	t := make(Trellis)
	for _, tag := range tags {
		t[tag] = make([]*Result, size)
		for i,_ := range t[tag] {
			t[tag][i] = &Result{}
		}
	}
	return &t
}

