package main

import(
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"regexp"
	"github.com/johnmcconnell/hmm"
)

type LabeledCache map[hmm.Tag]map[string]bool

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: train_file lexicon_file")
		os.Exit(-1)
	}
	cache := ParseLexicon(os.Args[2])
	tags := cache.Tags()
	i := hmm.UniformI(tags)
	t := hmm.UniformT(tags)
	e := hmm.UniformE(cache)

	train := ParseTraining(os.Args[1])

	em := hmm.NewEMLog(tags, train, i, t, e)

	i, t, e = EMLoop(2, em)
	lSentences := BuildLabeled(train, tags, i, t, e)

	PrintLabeledSentences(lSentences)
}

func EMLoop(index int, em *hmm.EMLog) (hmm.InitialState, hmm.Transition, hmm.Emission) {
	for x := 0; x < index; x++ {
		em = em.Next()
	  log.Printf("Finished %v EM training\n", x)
	}
	i := em.I()
	t := em.T()
	e := em.E()
	return i, t, e
}

func PrintLabeledSentences(sentences [][]hmm.LabeledWord) {
	for _, sentence := range sentences {
		PrintLabeledSentence(sentence)
	}
}

func PrintLabeledSentence(sentence []hmm.LabeledWord) {
	for _, word := range sentence {
		fmt.Sprintf("%s ", word)
	}
	fmt.Println()
}

func BuildLabeled(sentences [][]string,
tags []hmm.Tag, i hmm.InitialState,
t hmm.Transition, e hmm.Emission) [][]hmm.LabeledWord {
	labeledSentences := make([][]hmm.LabeledWord, len(sentences))
	for iS, sentence := range sentences {
	  v := hmm.NewViterbi(tags, sentence, &i, &t, &e)
		v.FillTrellis()
		labeled, err := v.Labeled()
		if err != nil {
			log.Println("Sentence unable to be labeled")
		  labeledSentences[iS] = make([]hmm.LabeledWord, 0)
		} else {
		  labeledSentences[iS] = labeled
		}
		if (iS % 100 == 0) {
			log.Println("Finished 100 sentences")
		}
	}
	return labeledSentences
}

// Tags ...
func (c *LabeledCache) Tags() []hmm.Tag {
	tags := make([]hmm.Tag, len(*c))
	i := 0
	for tag, _ := range *c {
		tags[i] = tag
		i += 1
	}
	return tags
}

// CheckTestParse ...
func CheckTestParse(cache LabeledCache, sentences [][]hmm.LabeledWord) {
	for _, sentence := range sentences {
		for _, word := range sentence {
			if !cache[word.Tag][word.Word] {
				fmt.Println("'%s' was not found in lexicon cache", word)
				os.Exit(-1)
			}
		}
	}
}

// ParseTraining ...
func ParseTraining(filename string) [][]string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	whiteSpace, err := regexp.Compile("\\s+")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	newLine, err := regexp.Compile("[\\s\\n]*\\n[\\s\\n]*")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	raw := string(bytes)
	lines := newLine.Split(raw, -1)
	if "" == lines[len(lines) - 1] {
	  lines = lines[:len(lines) - 1]
	}
	sentences := make([][]string, len(lines))
	for i, line := range lines {
		sentences[i] = whiteSpace.Split(line, -1)
	}
	return sentences
}

// ParseLexicon ...
func ParseLexicon(filename string) LabeledCache {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	newLine, err := regexp.Compile("[\\s\\n]*\\n[\\s\\n]*")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	raw := string(bytes)
	lines := newLine.Split(raw, -1)
	if "" == lines[len(lines) - 1] {
	  lines = lines[:len(lines) - 1]
	}
	lC := make(LabeledCache)
	for _, line := range lines {
		word, tags := ParseLexiconLine(line)
		for _, tag := range tags {
			if lC[hmm.Tag(tag)] == nil {
				lC[hmm.Tag(tag)] = make(map[string]bool)
			}
			lC[hmm.Tag(tag)][word] = true
		}
	}
	return lC
}

// ParseLexiconLine ...
func ParseLexiconLine(line string) (string, []string) {
	whiteSpace, err := regexp.Compile("\\s+")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
  parts := whiteSpace.Split(line, -1)
	word := parts[0]
	tags := parts[1:]
	return word, tags
}
