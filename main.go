package main

import (
	"fmt"
	"os"
	"git.enova.com/zsyed/utils"
)

type Config struct {
	Tags []Tag
	I InitialState
	T Transition
	E map[Tag]map[string]float64
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sequence config.yml")
		os.Exit(-1)
	}

	c := Config{}
	utils.ReadYAML(os.Args[2], &c)

	e := c.EmissionOfConfig()
	v := NewViterbi(c.Tags, os.Args[1], &c.I, &c.T, &e)
	v.FillTrellis()
	fmt.Printf("%s\n", v)
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
