package handlers

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Song struct {
	Title    string
	Subtitle string
	Link     string
}

var wg sync.WaitGroup
var mu sync.Mutex

func getTune(query string, allSongs *[]Song, wg *sync.WaitGroup) {

	c := colly.NewCollector(
		colly.AllowedDomains("https://get-tune.cc", "get-tune.cc"),
	)

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	// Lock the thread to prevent from being reused
	runtime.LockOSThread()

	// Lock the mutex to prevent the mutiple processes to write data at the same time
	// Strong Consistency
	mu.Lock()

	// Unlock the mutex
	defer mu.Unlock()

	defer wg.Done()

	defer runtime.UnlockOSThread()

	c.OnHTML(".playlist li", func(element *colly.HTMLElement) {

		link := element.Attr("data-mp3")
		songs := element.DOM

		song := songs.Find(".playlist-name").Find("b").Text()

		em := songs.Find(".playlist-name").Find("em").Text()

		if len(em) > 5 && !strings.Contains(strings.ToLower(em), "remix") && !strings.Contains(strings.ToLower(em), "mix") && !strings.Contains(strings.ToLower(em), "edit") && !strings.Contains(strings.ToLower(song), "mix") && !strings.Contains(strings.ToLower(song), "edit") && !strings.Contains(strings.ToLower(song), "remix") {
			*allSongs = append(*allSongs, Song{
				Title:    song,
				Subtitle: em,
				Link:     link,
			})
		}

	})
	err := c.Visit(fmt.Sprintf("https://get-tune.cc/search/f/%s/", strings.Join(strings.Split(query, " "), "+")))

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func Crawler(query string) []Song {

	songs := []Song{}

	wg.Add(1)
	go getTune(query, &songs, &wg)
	wg.Wait()

	return songs
}
