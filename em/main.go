package main

import(
	"os"
	"fmt"
	"io/ioutil"
	"github.com/johnmcconnell/hmm"
)

func Tain() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: train_file test_file config.yml")
		os.Exit(-1)
	}
	f := hmm.Forward{}
}

func ParseTraining(filename string) {
	ioutil.ReadFile(filename)
}
