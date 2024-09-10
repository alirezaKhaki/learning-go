package main

import (
	"fmt"
	"sync"

)

type SlowParser interface {
	Parse(string) string
}

// Simulate slow parser setup
func initParser() SlowParser {
	fmt.Println("Initializing parser...")
	return &MyParser{}
}

type MyParser struct{}

func (p *MyParser) Parse(input string) string {
	return "Parsed: " + input
}

// Global variables to keep track of the parser and sync.Once
var parser SlowParser
var once sync.Once

// Parse function that makes sure the parser is only initialized once
func Parse(dataToParse string) string {
	once.Do(func() {
		parser = initParser() // This code runs only once
	})
	return parser.Parse(dataToParse)
}

func main() {
	// Even though Parse is called twice, the parser is initialized only once
	fmt.Println(Parse("data1"))
	fmt.Println(Parse("data2"))
}
