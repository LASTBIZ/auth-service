package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"io"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/pkg/errors"
	"net/http"
	"net/url"
)

const FacebookUserAPI = "https://graph.facebook.com/me?access_token="

type Facebook struct {
	storage provider.Storage
	conf    oauth2.Config
}

func NewFacebookProvider(
	clientID,
	clientSecret,
	redirectURL string,
	storage provider.Storage) provider.Provider {
	var conf = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"public_profile", "email"},
		RedirectURL:  redirectURL,
		Endpoint:     endpoints.Facebook,
	}
	return &Facebook{
		conf:    conf,
		storage: storage,
	}
}

type FacebookUserResult struct {
	Email       string
	FamilyName  string
	GivenName   string
	Locale      string
	Name        string
	VerifyEmail bool
	Picture     string
}

func (f Facebook) GenerateOAuthToken(state string) (url string) {
	return f.conf.AuthCodeURL(state)
}

func (f Facebook) Callback(code string) (*oauth2.Token, error) {
	token, err := f.conf.Exchange(context.Background(), code)
	return token, err
}

func (f Facebook) GetUser(token *oauth2.Token) (interface{}, error) {
	client := f.conf.Client(context.Background(), token)
	response, err := client.Get(FacebookUserAPI + url.QueryEscape(token.AccessToken))
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

	userBody := &FacebookUserResult{
		Email:       GoogleUserRes["email"].(string),
		VerifyEmail: false,
		Name:        GoogleUserRes["name"].(string),
		GivenName:   GoogleUserRes["first_name"].(string),
		FamilyName:  GoogleUserRes["last_name"].(string),
		Picture:     GoogleUserRes["picture"].(string),
	}

	return userBody, err
}
