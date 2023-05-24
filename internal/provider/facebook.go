package provider

import (
	"auth-service/internal/conf"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"io"
	"net/http"
	"net/url"
)

const FacebookUserAPI = "https://graph.facebook.com/me?access_token="

type Facebook struct {
	conf oauth2.Config
}

func NewFacebookProvider(
	cfg *conf.Providers_Facebook,
	redirectURL string) Provider {
	var config = oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"public_profile", "email"},
		RedirectURL:  redirectURL,
		Endpoint:     endpoints.Facebook,
	}
	return &Facebook{
		conf: config,
	}
}

func (f Facebook) GenerateOAuthToken(state string) (url string) {
	return f.conf.AuthCodeURL(state)
}

func (f Facebook) Callback(code string) (*oauth2.Token, error) {
	token, err := f.conf.Exchange(context.Background(), code)
	return token, err
}

func (f Facebook) GetUser(token *oauth2.Token) (*User, error) {
	client := f.conf.Client(context.Background(), token)
	response, err := client.Get(FacebookUserAPI + url.QueryEscape(token.AccessToken))
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
		VerifyEmail: false,
		Name:        GoogleUserRes["name"].(string),
		GivenName:   GoogleUserRes["first_name"].(string),
		FamilyName:  GoogleUserRes["last_name"].(string),
		Picture:     GoogleUserRes["picture"].(string),
	}

	return userBody, err
}
