package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

const (
	Credentials = "credentials"
	BasicAuth   = "BasicAuth"
	BearerAuth  = "BearerAuth"
	ApiKeyAuth  = "ApiKeyAuth"
	OpenIDAuth  = "OpenIDAuth"
	OAuth2Auth  = "OAuth2Auth"
)

type ISecurity interface {
	Authorize(g *fiber.Ctx) error
	Provider() string
	Scheme() *openapi3.SecurityScheme
}

type Security struct {
	ISecurity
}
