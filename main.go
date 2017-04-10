package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const QCLONE_VERSION = "0.0.1"
const 

type Bundle struct {
	Name string 
	Repositories []Repository
}
type Repository struct {
	URL  string `json:"git_url"`
	Name string `json:"name"`
}

func checkErrorAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func getAllUserRepositories() ([]Repository, error) {
	var err error

	userToken := os.Getenv("QCLONE_GITHUB_TOKEN")
	if userToken == "" {
		return nil, errors.New("No GitHub token. Stopping")
	}

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	checkErrorAndPanic(err)
	req.Header.Add("User-Agent", "qClone "+QCLONE_VERSION)
	req.Header.Add("Authorization", "token "+userToken)

	repositories := make([]Repository, 0)
	for i := 1; ; i++ {
		targetUrl, err := url.Parse(fmt.Sprintf("https://api.github.com/user/repos?page=%d", i))
		checkErrorAndPanic(err)

		req.URL = targetUrl
		resp, err := httpClient.Do(req)
		checkErrorAndPanic(err)

		_repositories := make([]Repository, 0)
		jsonDecoder := json.NewDecoder(resp.Body)
		err = jsonDecoder.Decode(&_repositories)
		checkErrorAndPanic(err)
		if len(_repositories) == 0 {
			break
		}
		repositories = append(repositories, _repositories...)
	}
	return repositories, nil
}

func isBundle() bool {}
func isUserRepository(name string, repositories []Repository) bool {
	for _, r := range repositories {
		if name == r.Name {
			return true
		}
	}

	return false
}

func parseCommand() {
	arguments := os.Args[1:]
	repositories, err := getAllUserRepositories()
	checkErrorAndPanic(err)
	switch len(arguments) {
	case 0:
		for i, r := range repositories {
			fmt.Printf("[%d] %s\n", i+1, r.Name)
		}
	case 1:
		fmt.Println("Clone repository")
		if isUserRepository(arguments[0]) {
			fmt.Printf("I would clone: %s\n", arguments[0])
		}
	case 2:
		fmt.Printf("Arguments: %s,%s\n", arguments[0], arguments[1])
	}
}

func main() {
	parseCommand()
}
