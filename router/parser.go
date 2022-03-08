package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/long2ice/fibers/constants"
	"reflect"
	"strings"
	_ "unsafe"
)

//go:linkname ParseToStruct github.com/gofiber/fiber/v2.(*Ctx).parseToStruct
func ParseToStruct(ctx *fiber.Ctx, aliasTag string, out interface{}, data map[string][]string) error

//go:linkname equalFieldType github.com/gofiber/fiber/v2.equalFieldType
func equalFieldType(out interface{}, kind reflect.Kind, key string) bool

func headerParser(c *fiber.Ctx, model interface{}) error {
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
	if err := ParseToStruct(c, constants.HEADER, model, headerData); err != nil {
		return err
	}
	return nil
}
func paramsParser(c *fiber.Ctx, model interface{}) error {
	params := make(map[string][]string)
	for _, param := range c.Route().Params {
		params[param] = append(params[param], c.Params(param))
	}
	if err := ParseToStruct(c, constants.URI, model, params); err != nil {
		return err
	}
	return nil
}
