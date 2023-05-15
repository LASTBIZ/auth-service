package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/pkg/errors"
	"net/http"
	"net/url"
)

const GoogleUserAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type Google struct {
	storage provider.Storage
	conf    oauth2.Config
}

func NewGoogleProvider(
	clientID,
	clientSecret,
	redirectURL string,
	storage provider.Storage) provider.Provider {
	var conf = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"profile", "email"},
		RedirectURL:  redirectURL,
		Endpoint:     google.Endpoint,
	}
	return &Google{
		conf:    conf,
		storage: storage,
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

func (g Google) GetUser(token *oauth2.Token) (interface{}, error) {
	client := g.conf.Client(context.Background(), token)
	response, err := client.Get(GoogleUserAPI + url.QueryEscape(token.AccessToken))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
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

	userBody := &GoogleUserResult{
		Email:       GoogleUserRes["email"].(string),
		VerifyEmail: GoogleUserRes["verified_email"].(bool),
		Name:        GoogleUserRes["name"].(string),
		GivenName:   GoogleUserRes["given_name"].(string),
		FamilyName:  GoogleUserRes["family_name"].(string),
		Picture:     GoogleUserRes["picture"].(string),
		Locale:      GoogleUserRes["locale"].(string),
	}

	return userBody, err
}
