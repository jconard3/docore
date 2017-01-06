package main

import "os"
import "fmt"

import "golang.org/x/oauth2"
import "github.com/digitalocean/godo"

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func CreateClient() (godo.Client, error) {
	pat := os.Getenv("TOKEN")
	if pat == "" {
		fmt.Println("Environment variable 'TOKEN' not initialized.\nRun the following in your shell to initialize environment variables\nexport TOKEN=$(cat coreos_do_token)\n")
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	return *client, nil
}

func ListDroplets(client godo.Client) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	droplets, _, _ := client.Droplets.List(opt)
	for _, element := range droplets {
		fmt.Println(element.Name)
	}
}

func main() {
	client, _ := CreateClient()
	ListDroplets(client)
}
