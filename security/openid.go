package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type OpenID struct {
	Security
	ConnectUrl string
}

func (i *OpenID) Authorize(c *fiber.Ctx) error {
	return nil
}
func (i *OpenID) Provider() AuthType {
	return OpenIDAuth
}

func (i *OpenID) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:             "openIdConnect",
		OpenIdConnectUrl: i.ConnectUrl,
	}
}
