package main

import (
	"fmt"
	"os"
	"bytes"
	"git.enova.com/zsyed/utils"
)

type Config struct {
	Tags []Tag
	I InitialState
	T Transition
	E map[Tag]map[string]float64
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: algo sequence config.yml")
		os.Exit(-1)
	}

	c := Config{}
	utils.ReadYAML(os.Args[3], &c)

	seq := os.Args[2]
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
	} else {
		fmt.Printf("Undefined algo '%s'\n", algo)
		os.Exit(-1)
	}
}

func RunViterbi(sequence string, c Config) {
	e := c.EmissionOfConfig()
	v := NewViterbi(c.Tags, sequence, &c.I, &c.T, &e)
	v.FillTrellis()
	fmt.Printf("%s\n", v)
}

func RunForward(sequence string, c Config) {
	e := c.EmissionOfConfig()
	f := NewForward(c.Tags, sequence, &c.I, &c.T, &e)
	f.FillTrellis()
	fmt.Printf("%s\n", f)
}

func RunBackward(sequence string, c Config) {
	e := c.EmissionOfConfig()
	b := NewBackward(c.Tags, sequence, &c.I, &c.T, &e)
	b.FillTrellis()
	fmt.Printf("%s\n", b)
}


func RunGamma(sequence string, c Config) {
	e := c.EmissionOfConfig()
	g := NewGamma(c.Tags, sequence, &c.I, &c.T, &e)
	fmt.Printf("%s\n", g)
}

func RunCombined(sequence string, c Config) {
	e := c.EmissionOfConfig()
	b := NewBackward(c.Tags, sequence, &c.I, &c.T, &e)
	b.FillTrellis()
	f := NewForward(c.Tags, sequence, &c.I, &c.T, &e)
	f.FillTrellis()
	buffer := bytes.NewBufferString(fmt.Sprintf("Combined: '%s'\n|", sequence))
	for i, _ := range sequence {
		var sum float64 = 0.0
		for _, tag := range c.Tags {
			fP := (*f.trellis)[tag][i].probability
			bP := (*b.trellis)[tag][i].probability
			sum += (bP * fP)
		}
		buffer.WriteString(fmt.Sprintf(" i: '%v' %.4f |", i, sum))
	}
	fmt.Println(buffer.String())
}

func (c *Config) EmissionOfConfig() Emission {
	e := Emission{}
	for tag, probs  := range c.E {
		p := make(map[uint8]float64)
		for s, prob := range probs {
			p[s[0]] = prob
		}
		e[tag] = p
	}
	return e
}
