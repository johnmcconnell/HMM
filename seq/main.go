package main

import (
	"fmt"
	"os"
	"bytes"
	"strings"
	"git.enova.com/zsyed/utils"
	"github.com/johnmcconnell/hmm"
)

type Config struct {
	Tags []hmm.Tag
	I hmm.InitialState
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
	} else if algo == "combined" {
		RunCombined(seq, c)
	} else if algo == "gamma" {
		RunGamma(seq, c)
	} else if algo == "forwardlog" {
		RunForwardLog(seq, c)
	} else if algo == "backwardlog" {
		RunBackwardLog(seq, c)
	} else if algo == "gammalog" {
		RunGammaLog(seq, c)
		fmt.Printf("Undefined algo '%s'\n", algo)
		os.Exit(-1)
	}
}

func RunViterbi(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	v := hmm.NewViterbi(c.Tags, sequence, &c.I, &c.T, &e)
	v.FillTrellis()
	fmt.Printf("%s\n", v)
}

func RunForwardLog(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	f := hmm.NewForwardLog(c.Tags, sequence, &c.I, &c.T, &e)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunBackwardLog(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	f := hmm.NewBackwardLog(c.Tags, sequence, &c.I, &c.T, &e)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunGammaLog(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	f := hmm.NewGammaLog(c.Tags, sequence, &c.I, &c.T, &e)
	fmt.Printf("%s\n", f)
}

func RunForward(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	f := hmm.NewForward(c.Tags, sequence, &c.I, &c.T, &e)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunBackward(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	b := hmm.NewBackward(c.Tags, sequence, &c.I, &c.T, &e)
	b.FillTrellis()
	fmt.Printf("%s\n", b)
}

func RunGamma(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	g := hmm.NewGamma(c.Tags, sequence, &c.I, &c.T, &e)
	fmt.Printf("%s\n", g)
}

func RunCombined(sequence []string, c Config) {
	e := c.EmissionOfConfig()
	b := hmm.NewBackward(c.Tags, sequence, &c.I, &c.T, &e)
	b.FillTrellis()
	f := hmm.NewForward(c.Tags, sequence, &c.I, &c.T, &e)
	f.FillTrellis()
	buffer := bytes.NewBufferString(fmt.Sprintf("Combined: '%s'\n|", sequence))
	for i, _ := range sequence {
		var sum float64 = 0.0
		for _, tag := range c.Tags {
			fP := f.Result(tag, i).Probability
			bP := b.Result(tag, i).Probability
			sum += (bP * fP)
		}
		buffer.WriteString(fmt.Sprintf(" i: '%v' %.4f |", i, sum))
	}
	fmt.Println(buffer.String())
}

func (c *Config) EmissionOfConfig() hmm.Emission {
	e := hmm.Emission{}
	for tag, probs  := range c.E {
		p := make(map[string]float64)
		for s, prob := range probs {
			p[s] = prob
		}
		e[tag] = p
	}
	return e
}
