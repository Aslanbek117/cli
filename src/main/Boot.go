package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)
type UrlContext struct {
	ID 	int
	Url string
}

func httpGet(url string, client http.Client) map[string]int {

	var answer  = make(map[string]int)
	res, err := client.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	if res.StatusCode !=200 {
		answer[url] = -1
		log.Fatalln("Something wrong with provided URL:[",url, "], statusCode:[", res.StatusCode, "]" )
		return answer
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if nil != err {
		log.Fatalln(err)
	}

	answer[url] = strings.Count(string(body),pattern )
	return answer
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

//retrieve all values from channel
//make http Get request by url in channel
func getOccurrences(done <- chan bool, urlsChannel <- chan UrlContext, client http.Client) <- chan map[string]int {
	urlAndOccurrencesChannel := make(chan map[string]int)
	go func() {
		for object := range urlsChannel {
			select {
			case <- done:
				return
			case urlAndOccurrencesChannel <- httpGet(object.Url, client):
			}
		}
		close(urlAndOccurrencesChannel)
	}()
	return urlAndOccurrencesChannel

}
// merge all channels into one
// fan-in fan-out
func merge(done <- chan bool, channels ...<-chan map[string]int) <- chan map[string] int {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	faninChan := make(chan map[string]int)

	multiplex := func(someChannel <- chan map[string]int) {

		defer wg.Done()

		for finalAnswer := range someChannel{
			select {
			case <-done:
				return
			case faninChan <- finalAnswer:
			}

		}
	}

	for _, c := range channels {
		go multiplex(c)
	}

	go func() {

		wg.Wait()

		close(faninChan)
	}()

	return faninChan

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

	client := http.Client{Timeout: 30 * time.Second}

	done := make(chan bool)
	defer close(done)

	start := time.Now()
	urls := Produce(done, cmdArgs)

	workers := make([] <- chan map[string]int, len(cmdArgs))
	for i :=0; i< len(cmdArgs); i++ {
		workers[i] = getOccurrences(done, urls, client)
	}

	for n :=range merge(done, workers ...) {
		fmt.Println(n)
	}

	fmt.Printf("Took %fs to get occurrences from %d urls\n", time.Since(start).Seconds(), len(cmdArgs))

}