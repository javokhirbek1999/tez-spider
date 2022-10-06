package main

import (
	"fmt"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var wg sync.WaitGroup

var threadProfile = pprof.Lookup("threadcreate")

var mu sync.Mutex

func wiki2(c *colly.Collector, workerID int, workers []int) {

	runtime.LockOSThread()
	// mu.Lock()
	defer wg.Done()
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		e.ChildAttrs("a", "href")
		// for _, link := range links {
		// fmt.Println(link)
		// }
	})

	workers[0]++

	// mu.Unlock()

	runtime.UnlockOSThread()

	// done <- ""
}

func wiki(c *colly.Collector, done chan string, workerID int, workers []int) {

	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		e.ChildAttrs("a", "href")
		// fmt.Println(links)
		// print()
		// for _, link := range links {
		// fmt.Println(link)
		// }
	})

	workers[0]++

}

func init() {

	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {

	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)
	// done := make(chan string)

	start := time.Now()

	workers := []int{0}

	wg.Add(30000)
	for i := 0; i < 30000; i++ {
		// wg.Add(1)
		// go wiki(c, done, i, workers)
		go wiki2(c, i, workers)
	}

	c.Visit("https://en.wikipedia.org/wiki/Web_scraping")
	// <-done
	wg.Wait()

	fmt.Printf("%d workers are done in %v", workers[0], time.Since(start))
	fmt.Println("Threads created ", threadProfile.Count())
}
