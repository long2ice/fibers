package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/long2ice/fibers/constants"
	"reflect"
	"strings"
	_ "unsafe"
)

//go:linkname decoderBuilder github.com/gofiber/fiber/v2.decoderBuilder
func decoderBuilder(parserConfig fiber.ParserConfig) interface{}

//go:linkname equalFieldType github.com/gofiber/fiber/v2.equalFieldType
func equalFieldType(out interface{}, kind reflect.Kind, key string) bool

//go:linkname Decoder github.com/gofiber/fiber/v2/internal/schema.Decoder
type Decoder interface {
	Decode(dst interface{}, src map[string][]string) error
}

func ParseToStruct(aliasTag string, out interface{}, data map[string][]string) error {
	decoder := decoderBuilder(fiber.ParserConfig{
		SetAliasTag:       aliasTag,
		IgnoreUnknownKeys: true,
	})
	return decoder.(Decoder).Decode(out, data)
}
func HeaderParser(c *fiber.Ctx, model interface{}) error {
	headerData := make(map[string][]string)
	c.Request().Header.VisitAll(func(key, val []byte) {
		k := utils.UnsafeString(key)
		v := utils.UnsafeString(val)
		if strings.Contains(v, ",") && equalFieldType(model, reflect.Slice, k) {
			values := strings.Split(v, ",")
			for i := 0; i < len(values); i++ {
				headerData[k] = append(headerData[k], values[i])
			}
		} else {
			headerData[k] = append(headerData[k], v)
		}

	})
	return ParseToStruct(constants.HEADER, model, headerData)
}
func ParamsParser(c *fiber.Ctx, model interface{}) error {
	params := make(map[string][]string)
	for _, param := range c.Route().Params {
		params[param] = append(params[param], c.Params(param))
	}
	return ParseToStruct(constants.URI, model, params)
}
