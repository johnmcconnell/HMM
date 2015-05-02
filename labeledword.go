package hmm

import(
	"strings"
	"fmt"
	"os"
)

type LabeledWord struct {
	word string
	tag Tag
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
