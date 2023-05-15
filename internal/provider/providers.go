package provider

import "golang.org/x/oauth2"

type Provider interface {
	GenerateOAuthToken(state string) (url string)
	Callback(code string) (*oauth2.Token, error)
	GetUser(token *oauth2.Token) (interface{}, error)
}
