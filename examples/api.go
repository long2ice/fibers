package main

import (
	"fmt"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/security"
)

type TokenHeader struct {
	Token string `header:"token" validate:"required" json:"token" example:"test"`
}
type TestQuery struct {
	TokenHeader `embed:""`
	Name        string `query:"name" validate:"required" json:"name" description:"name of model" example:"test"`
	Enum        string `query:"enum" validate:"required,oneof=1 2" json:"enum" description:"enum of model" example:"1"`
	Optional    string `query:"optional" json:"optional"`
}

func (t *TestQuery) Handler(c *fiber.Ctx) error {
	user := c.Locals(security.Credentials).(security.User)
	fmt.Println(user)
	return c.JSON(t)
}

type TestQueryList struct {
	TokenHeader `embed:""`
	Name        string `query:"name" validate:"required" json:"name" description:"name of model" example:"test"`
}

func (t *TestQueryList) Handler(c *fiber.Ctx) error {
	user := c.Locals(security.Credentials).(security.User)
	fmt.Println(user)
	return c.JSON([]TestQueryList{*t})
}

type TestQueryPath struct {
	Name  string `query:"name" validate:"required" json:"name" description:"name of model" example:"test"`
	ID    int    `uri:"id" validate:"required" json:"id" description:"id of model" example:"1"`
	Token string `header:"token" validate:"required" json:"token" example:"test"`
}

func (t *TestQueryPath) Handler(c *fiber.Ctx) error {
	return c.JSON(t)
}

type TestForm struct {
	ID   int    `query:"id" validate:"required" json:"id" description:"id of model" example:"1"`
	Name string `form:"name" validate:"required" json:"name" description:"name of model" example:"test"`
	List []int  `form:"list" validate:"required" json:"list" description:"list of model"`
	Enum string `form:"enum" validate:"required,oneof=1 2" json:"enum" description:"enum of model" example:"1"`
}

func (t *TestForm) Handler(c *fiber.Ctx) error {
	fmt.Println(t)
	return c.JSON(t)
}

type TestJson struct {
	ID   int    `query:"id" validate:"required" json:"id" description:"id of model" example:"1"`
	Name string `json:"name" validate:"required"  description:"name of model" example:"test"`
	List []int  `json:"list" validate:"required"  description:"list of model"`
	Enum string `json:"enum" validate:"required,oneof=1 2"  description:"enum of model" example:"1"`
}

func (t *TestJson) Handler(c *fiber.Ctx) error {
	return c.JSON(t)
}

type TestNoModel struct {
	C string `cookie:"c" validate:"required" json:"cookie" description:"cookie is not supported in try it out of swagger ui" example:"test"`
}

func (t *TestNoModel) Handler(c *fiber.Ctx) error {
	return c.JSON(t)
}

type TestFile struct {
	File *multipart.FileHeader `form:"file" validate:"required" description:"file upload"`
}

func (t *TestFile) Handler(c *fiber.Ctx) error {
	fmt.Println(fiber.Map{"file": t.File.Filename})
	return c.JSON(fiber.Map{"file": t.File.Filename})
}
