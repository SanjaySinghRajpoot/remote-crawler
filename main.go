package main

import (
	"fmt"

	"github.com/SanjaySinghRajpoot/remote-crawler/config"
	"github.com/SanjaySinghRajpoot/remote-crawler/routes"
	"github.com/gin-gonic/gin"
	"gopkg.in/robfig/cron.v2"
)

var cnt = 10

// need a CRON Job for 24 hours set
func runCronJobs() {

	// c := colly.NewCollector()

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	// tempjobs := make([]models.Job, 0)

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	url := fmt.Sprintf("https://himalayas.app%s", link)

	// 	name := e.ChildText("h2.text-xl.font-medium.text-gray-900")
	// 	if name != "" && cnt != 0 {

	// 		saveJob := models.Job{
	// 			Name:        name,
	// 			Description: "",
	// 			URL:         url,
	// 		}

	// 		fmt.Print(saveJob)

	// 		tempjobs = append(tempjobs, saveJob)

	// 		cnt--
	// 	}
	// })

	// c.Visit("https://himalayas.app/jobs/developer")

	// for _, r := range tempjobs {

	// 	save := models.Job{
	// 		Name:        r.Name,
	// 		Description: r.Description,
	// 		URL:         r.URL,
	// 	}

	// 	result := config.DB.Create(&save)

	// 	if result.Error != nil {
	// 		fmt.Println("error")
	// 		return
	// 	}
	// }

	fmt.Println("himalayas website crawl")
}

func runHackerNews() {

	// c := colly.NewCollector()

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	// tempjobs := make([]models.Job, 0)

	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	url := fmt.Sprintf("https://himalayas.app%s", link)

	// 	name := e.ChildText("h2.text-xl.font-medium.text-gray-900")
	// 	if name != "" && cnt != 0 {

	// 		saveJob := models.Job{
	// 			Name:        name,
	// 			Description: "",
	// 			URL:         url,
	// 		}

	// 		fmt.Print(saveJob)

	// 		tempjobs = append(tempjobs, saveJob)

	// 		cnt--
	// 	}
	// })

	// c.Visit("https://himalayas.app/jobs/developer")

	// for _, r := range tempjobs {

	// 	save := models.Job{
	// 		Name:        r.Name,
	// 		Description: r.Description,
	// 		URL:         r.URL,
	// 	}

	// 	result := config.DB.Create(&save)

	// 	if result.Error != nil {
	// 		fmt.Println("error")
	// 		return
	// 	}
	// }

	fmt.Println("hacker news crawler")
}

func main() {

	config.Connect()

	cronJob := cron.New()

	cronJob.AddFunc("@every 1s", func() {
		// can use go routines
		runHackerNews()
	})

	cronJob.AddFunc("@every 1s", func() {
		runCronJobs()
	})

	cronJob.Start()

	fmt.Scanln()

	// starting the golang server
	router := gin.New()
	routes.UserRoute(router)
	router.Run(":8080")

}
