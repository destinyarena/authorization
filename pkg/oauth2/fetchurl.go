package oauth2

import (
	"fmt"
	"net/url"
)

func (o *oauth2) FetchURL(state *string) (string, error) {
	data := url.Values{}
	data.Set("response_type", o.Provider.ResponseType)
	data.Set("client_id", o.Provider.ClientID)

	if len(o.Provider.Scope) != 0 {
		data.Set("scope", o.Provider.Scope)
	}

	if o.Provider.RedirectURI != nil {
		data.Set("redirect_uri", *o.Provider.RedirectURI)
	}

	if state != nil {
		data.Set("state", *state)
	}

	params := data.Encode()

	return fmt.Sprintf("%s?%s", o.Provider.AuthorizeURL, params), nil
}
