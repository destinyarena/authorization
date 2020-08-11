package oauth2

type (
	// OauthPayload is the authorization redirect payload
	OauthPayload struct {
		Code string `json:"code" validate:"required"`
	}
)
