# Authorization

Authorization is responsible for handling Oauth2 for discord, bungie, and faceit and returning a signed JWT

## Endpoints

Most if not all endpoints in here for Oauth2 provider use

```
GET /api/v2/oauth/discord/url
GET /api/v2/oauth/discord/callback
GET /api/v2/oauth/bungie/url
GET /api/v2/oauth/bungie/callback
GET /api/v2/oauth/faceit/url
GET /api/v2/oauth/faceit/callback
```
