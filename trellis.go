package hmm

import(
	"fmt"
	"bytes"
	"math"
	"github.com/johnmcconnell/gologspace"
)

type Result struct {
	prevTag Tag
	Prob float64
}

type Trellis map[Tag][]*Result

// NewTrellis ...
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

// String ...
func (t *Trellis) String() string {
	return t.FormatString(gologspace.NumSpace{})
}

// FormatString ...
func (t *Trellis) FormatString(s gologspace.Space) string {
	buffer := bytes.NewBufferString("Trellis:\n")
	for tag, results := range *t {
		buffer.WriteString(t.RowString(tag, results, s))
	}
	return fmt.Sprintf(buffer.String())
}

// RowString ...
func (t *Trellis) RowString(tag Tag, results []*Result, s gologspace.Space) string {
	buffer := bytes.NewBufferString(fmt.Sprintf("|Tag: '%v'|", tag))
	for _, result := range results {
		buffer.WriteString(result.FormatString(s))
	}
	buffer.WriteString(fmt.Sprintf("\n"))
	return buffer.String()
}

// String ...
func (r *Result) FormatString(s gologspace.Space) string {
	p := s.Exit(r.Prob)
	return fmt.Sprintf("'%v': %.8fp|", r.prevTag, p)
}

// String ...
func (r *Result) String() string {
	return fmt.Sprintf("'%v': %.8fp|", r.prevTag, r.Prob)
}

// MaxRow ...
func (t *Trellis) MaxRow(tag Tag) int {
	mP, mI := math.Inf(-1), -1
	for i, r := range (*t)[tag] {
		if (mP <= r.Prob) {
			mP = r.Prob
			mI = i
		}
	}
	return mI
}

// MaxColumn ...
func (t *Trellis) MaxColumn(i int) Tag {
	mP, mTag := math.Inf(-1), Tag("")
	for tag, values := range (*t) {
		r := values[i]
		if (mP <= r.Prob) {
			mP = r.Prob
			mTag = tag
		}
	}
	return mTag
}
