package security

import (
	"encoding/base64"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Basic struct {
	Security
}

type User struct {
	Username string
	Password string
}

func decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
func (b *Basic) parseBasicAuth(c *fiber.Ctx) (User, error) {
	var user User
	auth := c.Get(fiber.HeaderAuthorization)
	if auth == "" {
		return user, fiber.NewError(fiber.StatusUnauthorized, "authorization header is missing")
	}
	if !strings.HasPrefix(strings.ToLower(auth), "basic ") {
		return user, fiber.NewError(fiber.StatusUnauthorized, "authorization header is not basic")
	}
	raw, err := decode(auth[6:])
	if err != nil {
		return user, fiber.ErrUnauthorized
	}
	credentials := strings.Split(string(raw), ":")
	user.Username = credentials[0]
	user.Password = credentials[1]
	return user, nil
}
func (b *Basic) Authorize(c *fiber.Ctx) error {
	user, err := b.parseBasicAuth(c)
	if err != nil {
		return err
	} else {
		b.Callback(c, user)
	}
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
