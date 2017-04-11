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
const QCLONE_DEFAULT_HOME = "~/git"

type Bundle struct {
	Name         string
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

func isBundle(name string) bool { return false }

func isUserRepository(name string, repositories []Repository) bool {
	for _, r := range repositories {
		if name == r.Name {
			return true
		}
	}

	return false
}

func usage() {
	fmt.Println("./qclone command [clone-to-path]")
}
func version() {
	fmt.Println("qClone Version " + QCLONE_VERSION)
}

func fetchCache() ([]Repository, error)           {}
func cacheRepositories(repositories []Repository) {}

func parseCommand() {
	var repositories []Repository
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		// We expect at least one aurgment
		usage()
		return
	}
	if arguments[0] == "version" {
		// We want to avoid fetching a repository
		// list for a non-repo related command
		version()
		return
	}

	// Now fetch repositories as commands beyond
	// this point will need them
	// TODO: cache this somewhere
	repositories, err := getAllUserRepositories()
	checkErrorAndPanic(err)

	switch len(arguments) {
	case 1:
		if arguments[0] == "list" {
			for i, r := range repositories {
				fmt.Printf("[%d] %s\n", i+1, r.Name)
			}
		}
		// Look for a bundle first
		if isBundle(arguments[0]) {
			if len(arguments) > 1 && arguments[1] != "" {
				fmt.Printf("I would clone bundle %s to %s\n", arguments[0], arguments[1])
			} else {
				fmt.Printf("I would clone bundle %s to %s\n", arguments[0], QCLONE_DEFAULT_HOME)
			}
			return
		}

		// Now look for a repository
		if isUserRepository(arguments[0], repositories) {
			if len(arguments) > 1 && arguments[1] != "" {
				fmt.Printf("I would clone %s to %s\n", arguments[0], arguments[1])
			} else {
				fmt.Printf("I would clone %s to %s\n", arguments[0], QCLONE_DEFAULT_HOME)
			}
			return
		}

		// TODO: look for a third-party user's repo
		// TODO: look for user/repo-name

		// Nothing to do at this point
		fmt.Println("Not sure what do at this point. Exiting")
	}
}

func main() {
	parseCommand()
}
