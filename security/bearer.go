package security

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type Bearer struct {
	Security
}

func (b *Bearer) Authorize(c *fiber.Ctx) error {
	auth := c.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "empty authentication")
	} else {
		splits := strings.Split(auth, "Bearer ")
		if len(splits) != 2 {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authentication string")
		} else {
			b.Callback(c, splits[1])
		}
		return c.Next()
	}
}

func (b *Bearer) Provider() AuthType {
	return BearerAuth
}

func (b *Bearer) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
}
