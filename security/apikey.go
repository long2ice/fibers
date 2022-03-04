package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type ApiKey struct {
	Security
	Name string
}

func (k *ApiKey) Authorize(c *fiber.Ctx) error {
	return c.Next()
}
func (k *ApiKey) Provider() string {
	return ApiKeyAuth
}

func (k *ApiKey) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type: "http",
		In:   "header",
		Name: k.Name,
	}
}
