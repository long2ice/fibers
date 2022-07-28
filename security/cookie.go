package security

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

type Cookie struct {
	Security
	Name string
}

func (k *Cookie) Authorize(c *fiber.Ctx) error {
	cookie := c.Cookies(k.Name)
	if cookie == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "empty cookie: "+k.Name)
	} else {
		k.Callback(c, cookie)
	}
	return c.Next()
}

func (k *Cookie) Provider() AuthType {
	return CookieAuth
}

func (k *Cookie) Scheme() *openapi3.SecurityScheme {
	return &openapi3.SecurityScheme{
		Type: "apiKey",
		In:   "cookie",
		Name: k.Name,
	}
}
