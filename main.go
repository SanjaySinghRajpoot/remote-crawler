package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

func getShortUrl(url string) string {

	baseUrl := "http://tinyurl.com/api-create.php?url="
	urlToShorten := url
	getReqUrl := baseUrl + urlToShorten

	response, err := http.Get(getReqUrl)
	if err != nil {
		log.Fatal(err)
	}

	// read response body
	body, error := io.ReadAll(response.Body)
	if error != nil {
		fmt.Println(error)
	}

	defer response.Body.Close()
	return string(body)
}

// need a CRON Job for 24 hours set
func runCronJobs() {

	// https://himalayas.app/jobs/developer

	baseURL := "www.himalayas.app"

	startingURL := "https://" + baseURL

	// url := []string{baseURL}

	c := colly.NewCollector(colly.AllowedDomains(), colly.Async(true))

	// Set the concurrency limit
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {

		log.Println("Something went wrong:", err)

	})

	c.OnResponse(func(r *colly.Response) {

		fmt.Println("Visited", r.Request.URL)

	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		e.Request.Visit(e.Attr("href"))

		// Print link
		url := fmt.Sprintf("https://himalayas.app%s", link)

		// get the shorted link
		shortedUrl := getShortUrl(url)

		name := e.ChildText("h2.text-xl.font-medium.text-gray-900")

		saveJob := models.Job{
			Name:        name,
			Description: "",
			URL:         shortedUrl,
		}

		result := config.DB.Where(models.Job{URL: shortedUrl}).FirstOrCreate(&saveJob)

		if result.Error != nil {
			fmt.Println("error")
			return
		}

	})

	fmt.Println("Starting crawl at: ", startingURL)

	if err := c.Visit(startingURL); err != nil {

		fmt.Println("Error on start of crawl: ", err)

	}
	c.Wait()
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

	fmt.Println("tweet was successfully")
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
// URL shortner - Done
// description shortener - Working on this
// check for valid link and description -> shorten them as per need
// Make a valid Tweet format that can be used for view level
// Add test cases to this project

func main() {

	config.Connect()

	cronJob := cron.New()

	cronJob.AddFunc("@every 1s", func() {
		runCronJobs()
	})

	cronJob.Start()

	fmt.Scanln()

	// getTweetFromDB()

	// starting the golang server
	router := gin.New()
	routes.UserRoute(router)
	router.Run(":8080")

}
