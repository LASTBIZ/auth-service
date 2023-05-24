package provider

import (
	"auth-service/internal/conf"
	"github.com/google/wire"
	"golang.org/x/oauth2"
)

// Set ProviderSet is provider providers.
// TODO edit config for provider
var Set = wire.NewSet(NewProviders)

type Provider interface {
	GenerateOAuthToken(state string) (url string)
	Callback(code string) (*oauth2.Token, error)
	GetUser(token *oauth2.Token) (*User, error)
}

type Struct struct {
	Providers map[string]Provider
}

func NewProviders(conf *conf.Providers) *Struct {
	providers := make(map[string]Provider, 0)
	providers["google"] = NewGoogleProvider(conf.Google, conf.RedirectUrl)
	providers["facebook"] = NewFacebookProvider(conf.Facebook, conf.RedirectUrl)
	return &Struct{
		Providers: providers,
	}
}

type User struct {
	Email       string
	FamilyName  string
	GivenName   string
	Name        string
	VerifyEmail bool
	Picture     string
}
