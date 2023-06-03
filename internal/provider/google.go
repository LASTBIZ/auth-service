package provider

import (
	"auth-service/internal/conf"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"net/url"
)

const GoogleUserAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type Google struct {
	conf oauth2.Config
}

func NewGoogleProvider(
	cfg *conf.Providers_Google,
	redirectURL string) Provider {
	var config = oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"profile", "email"},
		RedirectURL:  redirectURL + "/google",
		Endpoint:     google.Endpoint,
	}
	return &Google{
		conf: config,
	}
}

type GoogleUserResult struct {
	Email       string
	FamilyName  string
	GivenName   string
	Locale      string
	Name        string
	VerifyEmail bool
	Picture     string
}

func (g Google) GenerateOAuthToken(state string) (url string) {
	return g.conf.AuthCodeURL(state)
}

func (g Google) Callback(code string) (*oauth2.Token, error) {
	token, err := g.conf.Exchange(context.Background(), code)
	return token, err
}

func (g Google) GetUser(token *oauth2.Token) (*User, error) {
	client := g.conf.Client(context.Background(), token)
	response, err := client.Get(GoogleUserAPI + url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(500, "OAUTH_GET_USER_ERROR", "could not retrieve user")
	}
	var resBody bytes.Buffer
	defer response.Body.Close()
	_, err = io.Copy(&resBody, response.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &User{
		Email:       GoogleUserRes["email"].(string),
		VerifyEmail: GoogleUserRes["verified_email"].(bool),
		Name:        GoogleUserRes["name"].(string),
		GivenName:   GoogleUserRes["given_name"].(string),
		FamilyName:  GoogleUserRes["family_name"].(string),
		Picture:     GoogleUserRes["picture"].(string),
		//Locale:      GoogleUserRes["locale"].(string),
	}

	return userBody, err
}
