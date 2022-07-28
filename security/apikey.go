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
	auth := c.Get(k.Name)
	if auth == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "empty apikey")
	} else {
		k.Callback(c, auth)
	}
	return c.Next()
}

func (k *ApiKey) Provider() AuthType {
	return ApiKeyAuth
}

func (k *ApiKey) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type: "http",
		In:   "header",
		Name: k.Name,
	}
}
