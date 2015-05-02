package main

import(
	"os"
	"fmt"
	"io/ioutil"
	"regexp"
	"github.com/johnmcconnell/hmm"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: train_file test_file config.yml")
		os.Exit(-1)
	}
	ParseTraining(os.Args[1])
	ParseTest(os.Args[2])
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
