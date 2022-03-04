package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type Basic struct {
	Security
}

type User struct {
	Username string
	Password string
}

func (b *Basic) Authorize(c *fiber.Ctx) error {
	return c.Next()
}
func (b *Basic) Provider() string {
	return BasicAuth
}
func (b *Basic) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:   "http",
		Scheme: "basic",
	}
}
