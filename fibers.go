package fibers

import (
	"embed"
	"encoding/json"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/long2ice/fibers/router"
	"github.com/long2ice/fibers/swagger"
)

//go:embed templates/*
var templates embed.FS

type App struct {
	*fiber.App
	Swagger  *swagger.Swagger
	Routers  map[string]map[string]*router.Router
	subApps  map[string]*App
	rootPath string
}

func New(swagger *swagger.Swagger, config fiber.Config) *App {
	engine := html.NewFileSystem(http.FS(templates), ".html")
	config.Views = engine
	f := &App{App: fiber.New(config), Swagger: swagger, Routers: make(map[string]map[string]*router.Router), subApps: make(map[string]*App)}
	if swagger != nil {
		swagger.Routers = f.Routers
	}
	return f
}

func (g *App) Mount(path string, app *App) {
	app.rootPath = path
	app.App = g.App
	app.Swagger.Servers = append(app.Swagger.Servers, &openapi3.Server{
		URL: path,
	})
	g.subApps[path] = app
}

func (g *App) Group(path string, options ...Option) *Group {
	group := &Group{
		App:  g,
		Path: path,
	}
	for _, option := range options {
		option(group)
	}
	return group
}

func (g *App) Handle(path string, method string, r *router.Router) {
	r.Method = method
	r.Path = path
	if g.Routers[path] == nil {
		g.Routers[path] = make(map[string]*router.Router)
	}
	g.Routers[path][method] = r
}

func (g *App) Get(path string, router *router.Router) {
	g.Handle(path, fiber.MethodGet, router)
}

func (g *App) Post(path string, router *router.Router) {
	g.Handle(path, fiber.MethodPost, router)
}

func (g *App) Head(path string, router *router.Router) {
	g.Handle(path, fiber.MethodHead, router)
}

func (g *App) Patch(path string, router *router.Router) {
	g.Handle(path, fiber.MethodPatch, router)
}

func (g *App) Delete(path string, router *router.Router) {
	g.Handle(path, fiber.MethodDelete, router)
}

func (g *App) Put(path string, router *router.Router) {
	g.Handle(path, fiber.MethodPut, router)
}

func (g *App) Options(path string, router *router.Router) {
	g.Handle(path, fiber.MethodOptions, router)
}

func (g *App) init() {
	if g.Swagger == nil {
		return
	}
	g.App.Get(g.fullPath(g.Swagger.OpenAPIUrl), func(c *fiber.Ctx) error {
		return c.JSON(g.Swagger)
	})
	g.App.Get(g.fullPath(g.Swagger.DocsUrl), func(c *fiber.Ctx) error {
		options := `{}`
		if g.Swagger.SwaggerOptions != nil {
			data, err := json.Marshal(g.Swagger.SwaggerOptions)
			if err != nil {
				panic(err)
			}
			options = string(data)
		}
		return c.Render("templates/swagger", fiber.Map{
			"openapi_url":     g.fullPath(g.Swagger.OpenAPIUrl),
			"title":           g.Swagger.Title,
			"swagger_options": options,
		})
	})
	g.App.Get(g.fullPath(g.Swagger.RedocUrl), func(c *fiber.Ctx) error {
		options := `{}`
		if g.Swagger.RedocOptions != nil {
			data, err := json.Marshal(g.Swagger.RedocOptions)
			if err != nil {
				panic(err)
			}
			options = string(data)
		}
		return c.Render("templates/redoc", fiber.Map{
			"openapi_url":   g.fullPath(g.Swagger.OpenAPIUrl),
			"title":         g.Swagger.Title,
			"redoc_options": options,
		})
	})
	g.initRouters()
	g.Swagger.BuildOpenAPI()
}
func (g *App) initRouters() {
	for path, m := range g.Routers {
		path = g.fullPath(path)
		for method, r := range m {
			handlers := r.GetHandlers()
			if method == fiber.MethodGet {
				g.App.Get(path, handlers...)
			} else if method == fiber.MethodPost {
				g.App.Post(path, handlers...)
			} else if method == fiber.MethodHead {
				g.App.Head(path, handlers...)
			} else if method == fiber.MethodPatch {
				g.App.Patch(path, handlers...)
			} else if method == fiber.MethodDelete {
				g.App.Delete(path, handlers...)
			} else if method == fiber.MethodPut {
				g.App.Put(path, handlers...)
			} else if method == fiber.MethodOptions {
				g.App.Options(path, handlers...)
			} else {
				g.App.All(path, handlers...)
			}
		}
	}
}
func (g *App) fullPath(path string) string {
	return g.rootPath + path
}
func (g *App) Init() {
	g.init()
	for _, s := range g.subApps {
		s.init()
	}
}
func (g *App) Listen(addr string) error {
	g.Init()
	return g.App.Listen(addr)
}
