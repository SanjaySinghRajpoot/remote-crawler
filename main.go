package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/SanjaySinghRajpoot/remote-crawler/config"
	"github.com/SanjaySinghRajpoot/remote-crawler/models"
	"github.com/SanjaySinghRajpoot/remote-crawler/routes"
	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"gopkg.in/robfig/cron.v2"
)

var cnt = 10

func goGetEnv(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error, unable to load the .env file")
	}

	return os.Getenv(key)
}

// need a CRON Job for 24 hours set
func runCronJobs() {

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	tempjobs := make([]models.Job, 0)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		url := fmt.Sprintf("https://himalayas.app%s", link)

		name := e.ChildText("h2.text-xl.font-medium.text-gray-900")
		if name != "" && cnt != 0 {

			saveJob := models.Job{
				Name:        name,
				Description: "",
				URL:         url,
			}

			fmt.Print(saveJob)

			tempjobs = append(tempjobs, saveJob)

			cnt--
		}
	})

	c.Visit("https://himalayas.app/jobs/developer")

	for _, r := range tempjobs {

		save := models.Job{
			Name:        r.Name,
			Description: r.Description,
			URL:         r.URL,
		}

		result := config.DB.Create(&save)

		if result.Error != nil {
			fmt.Println("error")
			return
		}
	}

	fmt.Println("himalayas website crawl")
}

func sendTweeet(tweet string) {

	consumerKey := goGetEnv("ConsumerKey")
	consumerSecret := goGetEnv("ConsumerSecret")
	accessToken := goGetEnv("AccessToken")
	accessSecret := goGetEnv("AccessSecret")
	// prompt := os.Getenv("PROMPT")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		panic("Missing required environment variable")
	}

	fetched := tweet

	// From here on, Twitter POST API
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	// Necessary - Lambdas timing out.
	httpClient.Timeout = time.Second * 10

	path := "https://api.twitter.com/2/tweets"

	body := fmt.Sprintf(`{"text": "%s"}`, fetched)

	bodyReader := strings.NewReader(body)

	response, err := httpClient.Post(path, "application/json", bodyReader)

	if err != nil {
		log.Fatalf("Error when posting to twitter: %v", err)
	}

	defer response.Body.Close()
	log.Printf("Response was OK: %v", response)

	fmt.Println("tweet was succesfull")
}

func getTweetFromDB() {

	jobs := []models.Job{}
	config.DB.Limit(3).Find(&jobs)

	for _, job := range jobs {

		// tweet := fmt.Sprintf("Name: %s \n , Description: %s, \n link: %s", job.Name, job.Description, job.URL)

		tweet := fmt.Sprintf("Name: %s ", job.Name)

		sendTweeet(tweet)
	}

}

// Things to add ------------------------
// URL shortner
// description shortener
// check for valid link and description -> shorten them as per need
// Make a valid Tweet format that can be used for view level

func main() {

	config.Connect()

	cronJob := cron.New()

	cronJob.AddFunc("@every 1s", func() {
		runCronJobs()
	})

	cronJob.Start()

	fmt.Scanln()

	getTweetFromDB()

	// starting the golang server
	router := gin.New()
	routes.UserRoute(router)
	router.Run(":8080")

}
