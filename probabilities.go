package main

type InitialState map[string]float64
type Transition map[Tag]map[Tag]float64
type Emission map[Tag]map[string]float64
