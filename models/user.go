package models

type User struct {
	Name string `json:"name"`
}

type Job struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}
