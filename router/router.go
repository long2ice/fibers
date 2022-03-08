package router

import (
	"container/list"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/long2ice/fibers/security"
	"reflect"
)

type IAPI interface {
	Handler(context *fiber.Ctx) error
}
type Router struct {
	Handlers    *list.List
	Path        string
	Method      string
	Summary     string
	Description string
	Deprecated  bool
	ContentType string
	Tags        []string
	API         IAPI
	OperationID string
	Exclude     bool
	Securities  []security.ISecurity
	Response    Response
}

var validate = validator.New()

func BindModel(api IAPI) fiber.Handler {
	return func(c *fiber.Ctx) error {
		model := reflect.New(reflect.TypeOf(api).Elem()).Interface()
		if err := headerParser(c, model); err != nil {
			return err
		}
		if err := c.QueryParser(model); err != nil {
			return err
		}
		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut {
			if err := c.BodyParser(model); err != nil {
				return err
			}
		}
		if err := paramsParser(c, model); err != nil {
			return err
		}
		if err := validate.Struct(model); err != nil {
			return err
		}
		if err := copier.Copy(api, model); err != nil {
			return err
		}
		return c.Next()
	}
}

func (router *Router) GetHandlers() []fiber.Handler {
	var handlers []fiber.Handler
	for _, s := range router.Securities {
		handlers = append(handlers, s.Authorize)
	}
	for h := router.Handlers.Front(); h != nil; h = h.Next() {
		if f, ok := h.Value.(fiber.Handler); ok {
			handlers = append(handlers, f)
		}
	}
	handlers = append(handlers, router.API.Handler)
	return handlers
}

func New(api IAPI, options ...Option) *Router {
	r := &Router{
		Handlers: list.New(),
		API:      api,
		Response: make(Response),
	}
	r.Handlers.PushBack(BindModel(api))
	for _, option := range options {
		option(r)
	}
	return r
}
