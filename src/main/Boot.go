package main

import (
	"flag"
	"log"
	"os"
)
type UrlContext struct {
	ID 	int
	Url string
}

// at first need to populate a channel
// this function takes as parameter urls that contains all urls that provided by cli Args and produce them to channel
func Produce(done <-chan bool, urls[]string) <- chan UrlContext {
	urlsChannel := make(chan UrlContext)
	urlContext := make([]UrlContext, len(urls))

	for i :=0; i< len(urls);i++ {
		urlContext[i] = UrlContext{i, urls[i]}
	}

	go func() {
		for _, context := range urlContext {
			select {
			case <-done:
				return
			case urlsChannel <- context:
			}

		}
		close(urlsChannel)
	}()
	return urlsChannel
}

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