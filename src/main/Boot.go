package main

import (
	"flag"
	"log"
	"os"
)
type arrayFlags []string


func (i *arrayFlags) String() string {
	return "someVal"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var cmdArgs arrayFlags

const DEFAULT_PATTERN = "world"
var pattern string

func main() {
	flag.Var(&cmdArgs, "url", "url")
	flag.StringVar(&pattern, "pattern", DEFAULT_PATTERN, "search by provided pattern")

	flagset :=make(map[string]bool)

	flag.Visit(func(f * flag.Flag) { flagset[f.Name] = true})

	if !flagset["pattern"] {
		log.Println("pattern not explicitly set, using default:[", DEFAULT_PATTERN, "]")
	}

	if !flagset["url"] {
		log.Println("url not explicitly set, terminating...")
		os.Exit(1)
	}

}