package main

import(
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"regexp"
	"github.com/johnmcconnell/hmm"
)

type WordCache map[string]bool
type LabeledCache map[hmm.Tag]map[string]bool
type ConfusionMatrix map[hmm.Tag]map[hmm.Tag]float64

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: train_file test_file")
		os.Exit(-1)
	}
	train := ParseTest(os.Args[1])
	test := ParseTest(os.Args[1])

	CheckParse(train, test)
	a := Accuracy(train, test)
	m := ConfustionMatrix(train, test)
	fmt.Printf("Accuracy: %.4fp\n", a * 100)
}

func ConfusionMaxtrix(s1, s2 [][]hmm.LabeledWord) ConfusionMatrix {
	m := make(ConfustionMatrix)
	for iS, _ := range s1 {
		for iW, _ := range s1[iS] {
			l1, l2 := s1[iS][iW], s2[iS][iW]
			if (l1 != l2) {
				if (m[l1.Tag] == nil) {
					m[l1.Tag] = make(map[hmm.Tag]float64)
				}
				m[l1.Tag][l2.Tag] += 1.0
			}
		}
	}
  return m
}

func Accuracy(s1, s2 [][]hmm.LabeledWord) float64 {
	total := 0.0
	correct := 0.0
	for iS, _ := range s1 {
		for iW, _ := range s1[iS] {
			l1, l2 := s1[iS][iW], s2[iS][iW]
			if (l1 == l2) {
				correct += 1
			}
			total += 1.0
		}
	}
  return correct / total
}

// CheckParse ...
func CheckParse(s1, s2 [][]hmm.LabeledWord) {
	l1, l2 := len(s1), len(s2)
	if (l2 != l1) {
		log.Printf("s1=%v vs s2=%v", l1, l2)
		os.Exit(-1)
	}
	for i, _ := range s1 {
		l1, l2 := len(s1[i]), len(s2[i])
		if (l2 != l1) {
			log.Printf("s1[%v]=%v vs s2[%v]=%v", l1, l2)
			os.Exit(-1)
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

// ParseTest ...
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
