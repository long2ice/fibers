package router

import (
	"container/list"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/long2ice/fibers/security"
	"github.com/mcuadros/go-defaults"
)

type Model any

type Router struct {
	Handlers            *list.List
	Path                string
	Method              string
	Summary             string
	Description         string
	Deprecated          bool
	RequestContentType  string
	ResponseContentType string
	Tags                []string
	API                 fiber.Handler
	Model               Model
	OperationID         string
	Exclude             bool
	Securities          []security.ISecurity
	Response            Response
}

var validate = validator.New()

func BindModel(req interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		model := reflect.New(reflect.TypeOf(req).Elem()).Interface()
		if err := HeaderParser(c, model); err != nil {
			return err
		}
		if err := CookiesParser(c, model); err != nil {
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
		if err := ParamsParser(c, model); err != nil {
			return err
		}
		defaults.SetDefaults(model)
		if err := validate.Struct(model); err != nil {
			return err
		}
		if err := copier.Copy(req, model); err != nil {
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
	handlers = append(handlers, router.API)
	return handlers
}

func New[T Model, F func(c *fiber.Ctx, req T) error](f F, options ...Option) *Router {
	var model T
	h := BindModel(&model)
	r := &Router{
		Handlers: list.New(),
		Response: make(Response),
		API: func(ctx *fiber.Ctx) error {
			return f(ctx, model)
		},
		Model: model,
	}
	for _, option := range options {
		option(r)
	}

	r.Handlers.PushBack(h)
	return r
}

func (router *Router) WithSecurity(securities ...security.ISecurity) *Router {
	Security(securities...)(router)
	return router
}

func (router *Router) WithResponses(response Response) *Router {
	Responses(response)(router)
	return router
}

func (router *Router) WithHandlers(handlers ...fiber.Handler) *Router {
	Handlers(handlers...)(router)
	return router
}

func (router *Router) WithTags(tags ...string) *Router {
	Tags(tags...)(router)
	return router
}

func (router *Router) WithSummary(summary string) *Router {
	Summary(summary)(router)
	return router
}

func (router *Router) WithDescription(description string) *Router {
	Description(description)(router)
	return router
}

func (router *Router) WithDeprecated() *Router {
	Deprecated()(router)
	return router
}

func (router *Router) WithOperationID(ID string) *Router {
	OperationID(ID)(router)
	return router
}

func (router *Router) WithExclude() *Router {
	Exclude()(router)
	return router
}

func (router *Router) WithContentType(contentType string, contentTypeType ContentTypeType) *Router {
	ContentType(contentType, contentTypeType)(router)
	return router
}
