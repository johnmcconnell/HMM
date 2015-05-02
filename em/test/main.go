package main

import(
	"os"
	"fmt"
	"io/ioutil"
	"regexp"
	"github.com/johnmcconnell/hmm"
)

type LabeledCache map[hmm.Tag]map[string]bool

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: train_file test_file lexicon_file")
		os.Exit(-1)
	}
	cache := ParseLexicon(os.Args[3])
	// train := ParseTraining(os.Args[1])
	test := ParseTest(os.Args[2])
	CheckTestParse(cache, test)

	// fmt.Printf("%s\n", test)
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

func ParseTest(filename string) [][]hmm.LabeledWord {
	sentences := ParseTraining(filename)
	labeledSentences := make([][]hmm.LabeledWord, len(sentences))
	for iS, sentence := range sentences {
		labeledSentences[iS] = make([]hmm.LabeledWord, len(sentence))
		for iW, word := range sentence {
			labeledSentences[iS][iW] = hmm.ParseLabeledWord(word)
		}
	}
	return labeledSentences
}
