package oauth2

import "net/http"

type (
	oauth2 struct {
		Provider *Provider
		Client   *http.Client
	}

	// Oauth2 exports all the core functions used to interact with the oauth2 package
	Oauth2 interface {
		GetToken(string, interface{}) error
		FetchURL(*string) (string, error)
		GetBaseAPI() string
	}
)

// New returns a new oauth2 instance
func New(provider *Provider) (Oauth2, error) {
	return &oauth2{
		Provider: provider,
		Client:   new(http.Client),
	}, nil
}
