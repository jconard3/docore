package client

import (
	"context"
	"errors"

	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var Ctx = context.Background()

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
		err := errors.New("Could not read value \"do_token\" from config file.")
		return godo.Client{}, err
	}
	pat := viper.GetString("do_token")

	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	return *client, nil
}
