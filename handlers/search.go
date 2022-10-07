package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/gocolly/colly"
)

type Song struct {
	Title    string
	Subtitle string
	Link     string
}

func (s *Song) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)

	return encoder.Encode(s)
}

func Crawler(query string) []Song {

	c := colly.NewCollector()

	allSongs := []Song{}

	c.OnHTML(".playlist li", func(element *colly.HTMLElement) {

		link := element.Attr("data-mp3")
		songs := element.DOM

		song := songs.Find(".playlist-name").Find("b").Text()

		em := songs.Find(".playlist-name").Find("em").Text()

		if len(em) > 5 && !strings.Contains(strings.ToLower(em), "remix") && !strings.Contains(strings.ToLower(em), "mix") && !strings.Contains(strings.ToLower(em), "edit") && !strings.Contains(strings.ToLower(song), "mix") && !strings.Contains(strings.ToLower(song), "edit") && !strings.Contains(strings.ToLower(song), "remix") {
			allSongs = append(allSongs, Song{
				Title:    song,
				Subtitle: em,
				Link:     link,
			})
		}

	})

	err := c.Visit(fmt.Sprintf("https://get-tune.cc/search/f/%s/", strings.Join(strings.Split(query, " "), "+")))

	if err != nil {
		fmt.Printf("\nError: %v\n", err)
	}

	return allSongs
}
