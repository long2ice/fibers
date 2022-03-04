package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type Bearer struct {
	Security
}

func (b *Bearer) Authorize(c *fiber.Ctx) error {
	return c.Next()
}
func (b *Bearer) Provider() string {
	return BearerAuth
}

func (b *Bearer) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
}
