package main

import (
	"fmt"

	"gopkg.in/robfig/cron.v2"
)

var cnt = 10

func runCronJobs() {
	// 2
	s := cron.New()

	s.AddFunc("@every 1s", func() {
		fmt.Println("---cron---")
	})

	// 4
	s.Start()
}

func main() {

	// config.Connect()

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

	// // starting the golang server
	// router := gin.New()
	// routes.UserRoute(router)
	// router.Run(":8080")

	runCronJobs()
	fmt.Scanln()
}
