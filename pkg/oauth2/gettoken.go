package oauth2

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (o *oauth2) GetToken(code string, respStruct interface{}) error {
	base := fmt.Sprintf("%s/%s", o.Provider.BaseAPI, o.Provider.TokenURI)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("scope", o.Provider.Scope)

	if o.Provider.RedirectURI != nil {
		data.Set("redirect_uri", *o.Provider.RedirectURI)
	}

	req, err := http.NewRequest("POST", base, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	creds := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", o.Provider.ClientID, o.Provider.ClientSecret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := o.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("Error code: %d", resp.StatusCode)
	}

	json.Unmarshal(body, respStruct)

	return nil
}
