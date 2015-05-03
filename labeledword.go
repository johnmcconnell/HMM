package hmm

import(
	"strings"
	"fmt"
	"os"
)

type LabeledWord struct {
	Word string
	Tag Tag
}

func (w *LabeledWord) String() string {
	return fmt.Sprintf("%s_%s", w.Word, w.Tag)
}

func ParseLabeledWord(word string) LabeledWord {
	split := strings.Split(word, "_")
	if len(split) < 2 {
		fmt.Printf("'%s' is not labelable\n", word)
		os.Exit(-1)
	}
	literal := split[0]
	tag := Tag(split[1])
	return LabeledWord{literal, tag}
}
