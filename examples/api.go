package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/security"
	"mime/multipart"
)

type TokenHeader struct {
	Token string `header:"token" validate:"required" json:"token" default:"test"`
}
type TestQuery struct {
	TokenHeader `embed:""`
	Name        string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	Enum        string `query:"enum" validate:"required,oneof=1 2" json:"enum" description:"enum of model" default:"1"`
	Optional    string `query:"optional" json:"optional"`
}

func (t *TestQuery) Handler(c *fiber.Ctx) error {
	user := c.Locals(security.Credentials).(security.User)
	fmt.Println(user)
	return c.JSON(t)
}

type TestQueryList struct {
	TokenHeader `embed:""`
	Name        string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
}

func (t *TestQueryList) Handler(c *fiber.Ctx) error {
	user := c.Locals(security.Credentials).(security.User)
	fmt.Println(user)
	return c.JSON([]TestQueryList{*t})
}

type TestQueryPath struct {
	Name  string `query:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	ID    int    `uri:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Token string `header:"token" validate:"required" json:"token" default:"test"`
}

func (t *TestQueryPath) Handler(c *fiber.Ctx) error {
	return c.JSON(t)
}

type TestForm struct {
	ID   int    `query:"id" validate:"required" json:"id" description:"id of model" default:"1"`
	Name string `form:"name" validate:"required" json:"name" description:"name of model" default:"test"`
	List []int  `form:"list" validate:"required" json:"list" description:"list of model"`
	Enum string `form:"enum" validate:"required,oneof=1 2" json:"enum" description:"enum of model" default:"1"`
}

func (t *TestForm) Handler(c *fiber.Ctx) error {
	fmt.Println(t)
	return c.JSON(t)
}

type TestNoModel struct {
}

func (t *TestNoModel) Handler(c *fiber.Ctx) error {
	return c.JSON(nil)
}

type TestFile struct {
	File *multipart.FileHeader `form:"file" validate:"required" description:"file upload"`
}

func (t *TestFile) Handler(c *fiber.Ctx) error {
	fmt.Println(fiber.Map{"file": t.File.Filename})
	return c.JSON(fiber.Map{"file": t.File.Filename})
}
