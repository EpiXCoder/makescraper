package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type PostData struct {
	Title      string `json:"title"`
	Link       string `json:"link"`
	ScrapeTime string `json:"scrape_time"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.reddit.com", "reddit.com"),
	)

	var data []PostData
	var count int

	c.OnHTML(`a[id^="post-title-t3"]`, func(e *colly.HTMLElement) {
		if count >= 10 {
			return 
		}
		title := strings.TrimSpace(e.Text) 
		link := e.Request.AbsoluteURL(e.Attr("href"))
		scrapeTime := time.Now().Format(time.RFC3339) // Current time in RFC3339 format

		if title != "" && link != "" {
			data = append(data, PostData{
				Title:      title,
				Link:       link,
				ScrapeTime: scrapeTime,
			})
			count++
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.reddit.com/r/cybersecurity/top/?t=day")

	fmt.Println("Scraped data:")
	for _, d := range data {
		fmt.Printf("Title: %s, Link: %s, Scrape Time: %s\n", d.Title, d.Link, d.ScrapeTime)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return
	}
	fmt.Println("JSON data:", string(jsonData))

	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
	}
}

