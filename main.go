package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, _ := client.Repositories.List(ctx, "", nil)
	if len(repos) > 0 {
		for {
			fmt.Println("Select a repository to clone:")
			for i, repo := range repos {
				fmt.Println(i, ": ", *repo.Name)
			}

			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("#: ")
			selection, _ := reader.ReadString('\n')
			if selection == "q\n" {
				break
			}

			toNumClean := strings.Split(selection, "\n")[0]
			toNum, err := strconv.Atoi(toNumClean)
			if err != nil {
				panic(err)
			}
			if toNum < 0 || toNum > len(repos) {
				fmt.Println("Not a valid selection")
				continue
			} else {
				fmt.Println("You selected: ", *repos[toNum].Name, "(", *repos[toNum].URL, ")")

				_, err := git.PlainClone("/tmp/", false, &git.CloneOptions{
					URL:      *repos[toNum].CloneURL,
					Progress: os.Stdout,
				})

				if err != nil {
					panic(err)
				}
				break
			}
		}
	}
}
