package main

import (
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/security"
)

type TokenHeader struct {
	Token string `header:"token" validate:"required" json:"token" example:"test"`
}
type TestQueryReq struct {
	TokenHeader `       embed:""`
	Name        string `         query:"name"     validate:"required"           json:"name"     description:"name of model" example:"test"`
	Enum        string `         query:"enum"     validate:"required,oneof=1 2" json:"enum"     description:"enum of model" example:"1"`
	Optional    string `         query:"optional"                               json:"optional"`
}

func TestQuery(c *fiber.Ctx, req TestQueryReq) error {
	user := c.Locals(security.Credentials).(security.User)
	fmt.Println(user)
	return c.JSON(req)
}

type TestQueryListReq struct {
	TokenHeader `       embed:""`
	Name        string `         query:"name" validate:"required" json:"name" description:"name of model" example:"test"`
}

func TestQueryList(c *fiber.Ctx, req TestQueryListReq) error {
	user := c.Locals(security.Credentials).(security.User)
	fmt.Println(user)
	return c.JSON([]TestQueryListReq{req})
}

type TestQueryPathReq struct {
	Name  string    `query:"name" validate:"required" json:"name"  description:"name of model" example:"test"`
	ID    int       `             validate:"required" json:"id"    description:"id of model"   example:"1"    uri:"id"`
	UUID  uuid.UUID `query:"uuid" validate:"required" json:"uuid" description:"uuid of model"`
	Token string    `             validate:"required" json:"token"                             example:"test"          header:"token"`
	Num   *int      `query:"num"  json:"num"         example:"1"`
}

func TestQueryPath(c *fiber.Ctx, req TestQueryPathReq) error {
	return c.JSON(req)
}

type TestFormReq struct {
	ID   int    `query:"id" validate:"required"           json:"id"   description:"id of model"   example:"1"`
	Name string `           validate:"required"           json:"name" description:"name of model" example:"test" form:"name"`
	List []int  `           validate:"required"           json:"list" description:"list of model"                form:"list"`
	Enum string `           validate:"required,oneof=1 2" json:"enum" description:"enum of model" example:"1"    form:"enum"`
}

func TestForm(c *fiber.Ctx, req TestFormReq) error {
	fmt.Println(req)
	return c.JSON(req)
}

type TestJsonReq struct {
	ID   int    `query:"id" validate:"required"           json:"id"   description:"id of model"   example:"1"`
	Name string `           validate:"required"           json:"name" description:"name of model" example:"test"`
	List []int  `           validate:"required"           json:"list" description:"list of model"`
	Enum string `           validate:"required,oneof=1 2" json:"enum" description:"enum of model" example:"1"`
}

func TestJson(c *fiber.Ctx, req TestJsonReq) error {
	return c.JSON(req)
}

func TestNoModel(c *fiber.Ctx) error {
	return c.SendString("no model")
}

type TestFileReq struct {
	File *multipart.FileHeader `form:"file" validate:"required" description:"file upload"`
}

func TestFile(c *fiber.Ctx, req TestFileReq) error {
	fmt.Println(fiber.Map{"file": req.File.Filename})
	return c.JSON(fiber.Map{"file": req.File.Filename})
}
