package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type OAuth2 struct {
	Security
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

func (i *OAuth2) Authorize(c *fiber.Ctx) error {
	return nil
}

func (i *OAuth2) Provider() AuthType {
	return OAuth2Auth
}

func (i *OAuth2) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type: "oauth2",
		Flows: &openapi3.OAuthFlows{
			AuthorizationCode: &openapi3.OAuthFlow{
				AuthorizationURL: i.AuthorizationURL,
				TokenURL:         i.TokenURL,
				RefreshURL:       i.RefreshURL,
				Scopes:           i.Scopes,
			},
		},
	}
}
