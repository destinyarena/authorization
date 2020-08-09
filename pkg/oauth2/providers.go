package oauth2

// Provider contains config data about an oauth2 provider ie. Discord, Bungie, Faceit, Google
type Provider struct {
	AuthorizeURL string
	BaseAPI      string
	ClientID     string
	ClientSecret string
	ResponseType string
	Scope        string
	TokenURI     string
	RedirectURI  *string
}
