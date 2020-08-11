package oauth2

// GetBaseAPI returns provider's baseAPI url
func (o *oauth2) GetBaseAPI() string {
	return o.Provider.BaseAPI
}

// GetProvider returns oauth2 provider
func (o *oauth2) GetProvider() *Provider {
	return o.Provider
}
