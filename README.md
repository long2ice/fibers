# Fiber + Swagger = Fibers

[![deploy](https://github.com/long2ice/fibers/actions/workflows/deploy.yml/badge.svg)](https://github.com/long2ice/fibers/actions/workflows/deploy.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/long2ice/fibers.svg)](https://pkg.go.dev/github.com/long2ice/fibers)

## Introduction

`Fibers` is a web framework based on `Fiber` and `Swagger`, which wraps `Fiber` and provides built-in swagger api docs
and request model validation.

## Why I build this project?

Previous I have used [FastAPI](https://github.com/tiangolo/fastapi), which gives me a great experience in api docs
generation, because nobody like writing api docs.

Now I use `Fiber` but I can't found anything like that, I found [swag](https://github.com/swaggo/swag) but which write
docs with comment is so stupid. So there is `Fibers`.

## Installation

```shell
go get -u github.com/long2ice/fibers
```

## Online Demo

You can see online demo at <https://fibers.long2ice.io/docs> or <https://fibers.long2ice.io/redoc>.

![](https://raw.githubusercontent.com/long2ice/fibers/dev/images/docs.png)
![](https://raw.githubusercontent.com/long2ice/fibers/dev/images/redoc.png)

And you can reference all usage in [examples](https://github.com/long2ice/fibers/tree/dev/examples).

## Usage

### Build Swagger

Firstly, build a swagger object with basic information.

```go
package examples

import (
  "github.com/getkin/kin-openapi/openapi3"
  "github.com/long2ice/fibers/swagger"
)

func NewSwagger() *swagger.Swagger {
  return swagger.New("Fibers", "Swagger + Fiber = Fibers", "0.1.0",
    swagger.License(&openapi3.License{
      Name: "Apache License 2.0",
      URL:  "https://github.com/long2ice/fibers/blob/dev/LICENSE",
    }),
    swagger.Contact(&openapi3.Contact{
      Name:  "long2ice",
      URL:   "https://github.com/long2ice",
      Email: "long2ice@gmail.com",
    }),
    swagger.TermsOfService("https://github.com/long2ice"),
  )
}
```

### Write API

Then make api struct which implement `router.IAPI`.

```go
package examples

import "github.com/gofiber/fiber/v2"

type TestQuery struct {
  Name string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
}

func (t *TestQuery) Handler(c *fiber.Ctx) error {
  return c.JSON(t)
}
```

#### All supported tags

| name          | description                                                     |
|---------------|-----------------------------------------------------------------|
| `query`       | binding query param                                             |
| `form`        | binding body param                                              |
| `uri`         | binding path param                                              |
| `header`      | binding header param                                            |
| `validate`    | [validator](https://github.com/go-playground/validator) support |
| `description` | swagger docs param description                                  |
| `example`     | swagger docs param example                                      |
| `default`     | swagger docs param default value                                |
| `embed`       | embed struct params                                             |

Note that the attributes in `TestQuery`? `Fibers` will validate request and inject it automatically, then you can use it
in handler easily.

### Write Router

Then write router with some docs configuration and api.

```go
package examples

var query = router.New(
  &TestQuery{},
  router.Summary("Test Query"),
  router.Description("Test Query Model"),
  router.Tags("Test"),
)
```

### Security

If you want to project your api with a security policy, you can use security, also they will be shown in swagger docs.

Current there is five kinds of security policies.

- `Basic`
- `Bearer`
- `ApiKey`
- `OpenID`
- `OAuth2`

```go
package main

var query = router.New(
  &TestQuery{},
  router.Summary("Test query"),
  router.Description("Test query model"),
  router.Security(&security.Basic{}),
)
```

Then you can get the authentication string by `c.Locals(security.Credentials)` depending on your auth type.

```go
package main

import "github.com/gofiber/fiber/v2"

func (t *TestQuery) Handler(c *fiber.Ctx) error {
  user := c.Locals(security.Credentials).(security.User)
  fmt.Println(user)
  return c.JSON(t)
}
```

### Mount Router

Then you can mount router in your application or group.

```go
package main

import "github.com/gofiber/fiber/v2"

func main() {
  app := fibers.New(NewSwagger(), fiber.Config{})
  queryGroup := app.Group("/query", fibers.Tags("Query"))
  queryGroup.Get("", query)
  queryGroup.Get("/:id", queryPath)
  queryGroup.Delete("", query)
  app.Get("/noModel", noModel)
}

```

### Start APP

Finally, start the application with routes defined.

```go
package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gofiber/fiber/v2"
  "github.com/long2ice/fibers"
)

func main() {
  app := fibers.New(NewSwagger(), fiber.Config{})
  app.Use(
    logger.New(),
    recover.New(),
    cors.New(),
  )
  subApp := fibers.New(NewSwagger(), fiber.Config{})
  subApp.Get("/noModel", noModel)
  app.Mount("/sub", subApp)
  app.Use(cors.New(cors.Config{
    AllowOrigins:     "*",
    AllowMethods:     "*",
    AllowHeaders:     "*",
    AllowCredentials: true,
  }))
  queryGroup := app.Group("/query", fibers.Tags("Query"))
  queryGroup.Get("/list", queryList)
  queryGroup.Get("/:id", queryPath)
  queryGroup.Delete("", query)

  app.Get("/noModel", noModel)

  formGroup := app.Group("/form", fibers.Tags("Form"), fibers.Security(&security.Bearer{}))
  formGroup.Post("/encoded", formEncode)
  formGroup.Put("", body)
  formGroup.Post("/file", file)

  log.Fatal(app.Listen(":8080"))
}
```

That's all! Now you can visit <http://127.0.0.1:8080/docs> or <http://127.0.0.1:8080/redoc> to see the api docs. Have
fun!

### Disable Docs

In some cases you may want to disable docs such as in production, just put `nil` to `fibers.New`.

```go
app = fibers.New(nil, fiber.Config{})
```

### SubAPP Mount

If you want to use sub application, you can mount another `SwaGin` instance to main application, and their swagger docs
is also separate.

```go
package main

import "github.com/gofiber/fiber/v2"

func main() {
  app := fibers.New(NewSwagger(), fiber.Config{})
  subApp := fibers.New(NewSwagger(), fiber.Config{})
  subApp.Get("/noModel", noModel)
  app.Mount("/sub", subApp)
}

```

## ThanksTo

- [kin-openapi](https://github.com/getkin/kin-openapi), OpenAPI 3.0 implementation for Go (parsing, converting,
  validation, and more).
- [Fiber](https://github.com/gofiber/fiber), Express inspired web framework written in Go.

## License

This project is licensed under the
[Apache-2.0](https://github.com/long2ice/fibers/blob/master/LICENSE)
License.
