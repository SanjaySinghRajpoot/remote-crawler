package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {

	// scrapeUrl := "https://himalayas.app/jobs"

	// c := colly.NewCollector(colly.AllowedDomains("https://himalayas.app/jobs", "himalayas.app/jobs"))

	// fmt.Println("-----w----")
	// // Find and visit all links
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	// c.OnHTML("span.sr-only", func(e *colly.HTMLElement) {
	// 	fmt.Println(e.Text)
	// })

	// c.Visit(scrapeUrl)

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Find and visit all links
	// c.OnHTML("a", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })

	c.OnHTML("h2.text-xl font-medium text-gray-900", func(h *colly.HTMLElement) {
		fmt.Println(h.Text)
	})

	c.Visit("https://himalayas.app/jobs")
}
