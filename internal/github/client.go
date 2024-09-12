package githubconsumer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Client struct {
	httpClient  *http.Client
	accessToken string
}

func New() *Client {
	c := Client{
		httpClient:  &http.Client{},
		accessToken: os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"),
	}

	return &c
}

func (c *Client) Octocat() {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/octocat", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"))
	req.Header.Set("X-GitHub-API-Version", "2022-11-28")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", bodyText)
}

func (c *Client) PullRequests(owner, repo string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+owner+"/"+repo+"/"+"pulls", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"))
	req.Header.Set("X-GitHub-API-Version", "2022-11-28")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", bodyText)
}
