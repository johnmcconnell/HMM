package main

import(
	"os"
	"fmt"
	"log"
	"math"
	"io/ioutil"
	"regexp"
	"github.com/johnmcconnell/hmm"
)

type WordCache map[string]bool
type LabeledCache map[hmm.Tag]map[string]bool

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: train_file lexicon_file")
		os.Exit(-1)
	}
	tagCache, wordCache := ParseLexicon(os.Args[2])
	tags := tagCache.Tags()
	words := wordCache.Words()

	log.Printf("len: %v", len(tags))
	i := hmm.UniformI(tags)
	t := hmm.UniformT(tags)
	e := hmm.UniformE2(tags, wordCache)

	train := ParseTraining(os.Args[1])

	em := hmm.NewEMLog2(tags, words, train, i, t, e)

	i, t, e = EMLoop(8, em)
	lSentences := BuildLabeled(train, tags, i, t, e)

	PrintLabeledSentences(lSentences)
}

func EMLoop(index int, em *hmm.EMLog) (hmm.InitialState, hmm.Transition, hmm.Emission) {
	i := em.I()
	t := em.T()
	e := em.E()
	for x := 0; x < index; x++ {
		em = em.Next()
	  log.Printf("Finished %v EM training\n", x)
		prevT := t
		prevE := e
		i = em.I()
		t = em.T()
		e = em.E()
		if Converges(prevT, t, prevE, e) {
	    log.Printf("Converges!\n", x)
			break
		} else {
	    log.Printf("Does not converge\n", x)
		}
	}
	return i, t, e
}

func Converges(t1, t2 hmm.Transition, e1, e2 hmm.Emission) bool {
	con := true
	for givenTag, _ := range t1 {
		for tag, _ := range t1[givenTag] {
			t1P := t1[givenTag][tag]
			t2P := t2[givenTag][tag]
			diff := math.Abs(t1P - t2P)
			if (diff > 0.0000001) {
				con = false
			}
			if (diff > 0.3) {
				log.Printf("diff '%v', T(%s|%s)", diff, tag, givenTag)
			}
		}
		for word, _ := range e2[givenTag] {
			e1P := e1[givenTag][word]
			e2P := e2[givenTag][word]
			diff := math.Abs(e1P - e2P)
			if (diff > 0.0000001) {
				con = false
			}
			if (diff > 0.3) {
				log.Printf("diff '%v', E(%s|%s)", diff, word, givenTag)
			}
		}
	}
	return con
}

func PrintLabeledSentences(sentences [][]hmm.LabeledWord) {
	for _, sentence := range sentences {
		PrintLabeledSentence(sentence)
	}
}

func PrintLabeledSentence(sentence []hmm.LabeledWord) {
	for _, word := range sentence {
		fmt.Printf("%s ", word.String())
	}
	fmt.Println()
}

func BuildLabeled(sentences [][]string,
tags []hmm.Tag, i hmm.InitialState,
t hmm.Transition, e hmm.Emission) [][]hmm.LabeledWord {
	labeledSentences := make([][]hmm.LabeledWord, len(sentences))
	for iS, sentence := range sentences {
	  v := hmm.NewViterbiLog(tags, sentence, &i, &t, &e)
		v.FillTrellis()
		labeled, err := v.Labeled()
		if err != nil {
			log.Printf("Failed Sentence: %v\n", len(sentence))
			log.Printf("Error: %v\n", err.Error())
		  labeledSentences[iS] = make([]hmm.LabeledWord, 0)
		} else {
		  labeledSentences[iS] = labeled
		}
		if (iS % 3000 == 0) {
			log.Println("Finished 3000 sentences")
		}
	}
	return labeledSentences
}

// Words ...
func (c *WordCache) Words() []string {
	words := make([]string, len(*c))
	i := 0
	for word, _ := range *c {
		words[i] = word
		i += 1
	}
	return words
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
func CheckParse(cache LabeledCache, sentences [][]hmm.LabeledWord) {
	for _, sentence := range sentences {
		for _, word := range sentence {
			if !cache[word.Tag][word.Word] {
				log.Printf("'%s' was not found in lexicon cache \n", word)
				log.Printf("C[%s] has [%s]\n", word.Tag, cache[word.Tag])
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
func ParseLexicon(filename string) (LabeledCache, WordCache) {
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
	wC := make(WordCache)
	for _, line := range lines {
		word, tags := ParseLexiconLine(line)
		for _, tag := range tags {
			if (tag == "") {
				log.Printf("Tag: '%s' is in [%s]\n", tag, tags)
				os.Exit(-1)
				continue
			}
			if (word == "") {
				log.Printf("Word: '%s' is blank in [%s]\n", word, line)
				os.Exit(-1)
				continue
			}
			t := hmm.Tag(tag)
			if lC[t] == nil {
				lC[t] = make(map[string]bool)
			}
			lC[t][word] = true
			wC[word] = true
		}
	}
	return lC, wC
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
