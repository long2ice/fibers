package fibers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/router"
	"github.com/long2ice/fibers/security"
)

type Group struct {
	*App
	Path       string
	Tags       []string
	Handlers   []fiber.Handler
	Securities []security.ISecurity
}
type Option func(*Group)

func Handlers(handlers ...fiber.Handler) Option {
	return func(g *Group) {
		for _, handler := range handlers {
			g.Handlers = append(g.Handlers, handler)
		}
	}
}

func Tags(tags ...string) Option {
	return func(g *Group) {
		if g.Tags == nil {
			g.Tags = tags
		} else {
			g.Tags = append(g.Tags, tags...)
		}
	}
}

func Security(securities ...security.ISecurity) Option {
	return func(g *Group) {
		for _, s := range securities {
			g.Securities = append(g.Securities, s)
		}
	}
}

func (g *Group) Handle(path string, method string, r *router.Router) {
	router.Handlers(g.Handlers...)(r)
	router.Tags(g.Tags...)(r)
	router.Security(g.Securities...)(r)
	g.App.Handle(g.Path+path, method, r)
}

func (g *Group) Get(path string, router *router.Router) {
	g.Handle(path, http.MethodGet, router)
}

func (g *Group) Post(path string, router *router.Router) {
	g.Handle(path, http.MethodPost, router)
}

func (g *Group) Head(path string, router *router.Router) {
	g.Handle(path, http.MethodHead, router)
}

func (g *Group) Patch(path string, router *router.Router) {
	g.Handle(path, http.MethodPatch, router)
}

func (g *Group) Delete(path string, router *router.Router) {
	g.Handle(path, http.MethodDelete, router)
}

func (g *Group) Put(path string, router *router.Router) {
	g.Handle(path, http.MethodPut, router)
}

func (g *Group) Options(path string, router *router.Router) {
	g.Handle(path, http.MethodOptions, router)
}

func (g *Group) Group(path string, options ...Option) *Group {
	group := &Group{
		App:        g.App,
		Path:       g.Path + path,
		Tags:       g.Tags,
		Handlers:   g.Handlers,
		Securities: g.Securities,
	}
	for _, option := range options {
		option(group)
	}
	return group
}
