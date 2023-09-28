package main

import (
	"encoding/json"
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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

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

func makeTweetUsingGPT(tweet models.Job) string {

	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{
		"model": "gpt-3.5-turbo",
		"messages": [
		  {
			"role": "system",
			"content": "You will be provided with a job title, and your task is to create a attractive tweet using the job title with 270 characters limit"
		  },
		  {
			"role": "user",
			"content": "%s"
		  }
		],
		"temperature": 0,
		"max_tokens": 256,
		"top_p": 1,
		"frequency_penalty": 0,
		"presence_penalty": 0
     }`, tweet.Name))
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", data)
	if err != nil {
		log.Fatal(err)
	}

	openAIKey := goGetEnv("OPEN_AI_KEY")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIKey)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObj ChatCompletionResponse

	unmarshalError := json.Unmarshal([]byte(bodyText), &responseObj)
	if unmarshalError != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return responseObj.Choices[0].Message.Content
}

func getTweetFromDB() {

	jobs := []models.Job{}
	config.DB.Limit(3).Find(&jobs)

	for _, job := range jobs {

		tweet := makeTweetUsingGPT(job)

		sendTweeet(tweet)
	}

}

func main() {

	config.Connect()

	cronJob := cron.New()

	// Cron Job set up to run on Weekly basis
	cronJob.AddFunc("@weekly", func() {
		runCronJobs()
	})

	cronJob.Start()

	// getTweetFromDB()

	// starting the golang server
	router := gin.New()
	routes.UserRoute(router)
	router.Run(":8080")

}
