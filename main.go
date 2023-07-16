package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Job struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// var tempJob []Job

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		url := fmt.Sprintf("https://himalayas.app%s", link)

		name := e.ChildText("h2.text-xl.font-medium.text-gray-900")
		if name != "" {

			saveJob := Job{
				Name:        name,
				Description: "",
				URL:         url,
			}

			fmt.Print(saveJob)
			fmt.Println("")
		}
	})

	// on Pagination handle this function

	c.Visit("https://himalayas.app/jobs/developer")
}
