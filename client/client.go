package client

import "os"
import "fmt"

import "golang.org/x/oauth2"
import "github.com/digitalocean/godo"
import "github.com/spf13/viper"

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
	if !viper.IsSet("do_token") {
		fmt.Println("Could not read value \"do_token\" from config file.")
		os.Exit(1)
	}
	pat := viper.GetString("do_token")

	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	return *client, nil
}
