package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/security"
)

type Option func(router *Router)

func Security(securities ...security.ISecurity) Option {
	return func(router *Router) {
		for _, s := range securities {
			router.Securities = append(router.Securities, s)
		}
	}
}

func Responses(response Response) Option {
	return func(router *Router) {
		router.Response = response
	}
}

func Handlers(handlers ...fiber.Handler) Option {
	return func(router *Router) {
		for _, handler := range handlers {
			router.Handlers.PushBack(handler)
		}
	}
}

func Tags(tags ...string) Option {
	return func(router *Router) {
		if router.Tags == nil {
			router.Tags = tags
		} else {
			router.Tags = append(router.Tags, tags...)
		}
	}
}

func Summary(summary string) Option {
	return func(router *Router) {
		router.Summary = summary
	}
}

func Description(description string) Option {
	return func(router *Router) {
		router.Description = description
	}
}

// Deprecated mark api is deprecated
func Deprecated() Option {
	return func(router *Router) {
		router.Deprecated = true
	}
}

func OperationID(ID string) Option {
	return func(router *Router) {
		router.OperationID = ID
	}
}

// Exclude exclude in docs
func Exclude() Option {
	return func(router *Router) {
		router.Exclude = true
	}
}

// ContentType Set request contentType
func ContentType(contentType string, contentTypeType ContentTypeType) Option {
	return func(router *Router) {
		if contentTypeType == ContentTypeRequest {
			router.RequestContentType = contentType
		} else {
			router.ResponseContentType = contentType
		}
	}
}
