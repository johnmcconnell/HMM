package main

import (
	"fmt"
	"os"
	"strings"
	"git.enova.com/zsyed/utils"
	"github.com/johnmcconnell/hmm"
	"github.com/johnmcconnell/gologspace"
)

type Config struct {
	Tags []hmm.Tag
	I hmm.Initial
	T hmm.Transition
	E map[hmm.Tag]map[string]float64
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: algo sequence config.yml")
		os.Exit(-1)
	}

	c := Config{}
	utils.ReadYAML(os.Args[3], &c)

	seq := strings.Split(os.Args[2], "")
	algo := os.Args[1]
	if algo == "viterbi" {
		RunViterbi(seq, c)
	} else if algo == "forward" {
		RunForward(seq, c)
	} else if algo == "backward" {
		RunBackward(seq, c)
	} else if algo == "gamma" {
		RunGamma(seq, c)
	} else if algo == "viterbilog" {
		RunViterbiLog(seq, c)
	} else if algo == "forwardlog" {
		RunForwardLog(seq, c)
	} else if algo == "backwardlog" {
		RunBackwardLog(seq, c)
	} else if algo == "gammalog" {
		RunGammaLog(seq, c)
	} else {
		fmt.Printf("Undefined algo '%s'\n", algo)
		os.Exit(-1)
	}
}

func RunViterbi(sequence []string, c Config) {
	s := gologspace.NumSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	v := hmm.NewViterbi(c.Tags, sequence, s, &I, &T, &E)
	v.FillTrellis()
	fmt.Printf("%s\n", v)
}

func RunViterbiLog(sequence []string, c Config) {
	s := gologspace.LogSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	v := hmm.NewViterbi(c.Tags, sequence, s, &I, &T, &E)
	v.FillTrellis()
	fmt.Printf("%s\n", v)
}

func RunForward(sequence []string, c Config) {
	s := gologspace.NumSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	f := hmm.NewForward(c.Tags, sequence, s, &I, &T, &E)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunForwardLog(sequence []string, c Config) {
	s := gologspace.LogSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	f := hmm.NewForward(c.Tags, sequence, s, &I, &T, &E)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunBackward(sequence []string, c Config) {
	s := gologspace.NumSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	b := hmm.NewBackward(c.Tags, sequence, s, &I, &T, &E)
	b.FillTrellis()
	fmt.Printf("%s\n", b)
}

func RunBackwardLog(sequence []string, c Config) {
	s := gologspace.LogSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	f := hmm.NewBackward(c.Tags, sequence, s, &I, &T, &E)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunGamma(sequence []string, c Config) {
	s := gologspace.NumSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	b := hmm.NewGamma(c.Tags, sequence, s, &I, &T, &E)
	fmt.Printf("%s\n", b)
}

func RunGammaLog(sequence []string, c Config) {
	s := gologspace.LogSpace{}
	E := c.Emission(s)
	T := c.Transition(s)
	I := c.Initial(s)
	f := hmm.NewGamma(c.Tags, sequence, s, &I, &T, &E)
	fmt.Printf("%s\n", f)
}

func (c *Config) Transition(s gologspace.Space) hmm.Transition {
	e := hmm.Transition{}
	for tag1, probs  := range c.T {
		p := make(map[hmm.Tag]float64)
		for tag2, prob := range probs {
			p[tag2] = s.Enter(prob)
		}
		e[tag1] = p
	}
	return e
}

func (c *Config) Emission(s gologspace.Space) hmm.Emission {
	e := hmm.Emission{}
	for tag, probs  := range c.E {
		p := make(map[string]float64)
		for word, prob := range probs {
			p[word] = s.Enter(prob)
		}
		e[tag] = p
	}
	return e
}

func (c *Config) Initial(s gologspace.Space) hmm.Initial {
	e := hmm.Initial{}
	for tag, p  := range c.I {
		e[tag] = s.Enter(p)
	}
	return e
}
