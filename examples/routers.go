package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/router"
	"github.com/long2ice/fibers/security"
)

var (
	query = router.New(
		TestQuery,
		router.Summary("Test query"),
		router.Description("Test query model"),
		router.Security(&security.Basic{}),
		router.Responses(router.Response{
			"200": router.ResponseItem{
				Model:       TestQueryReq{},
				Description: "response model description",
			},
		}),
	)

	queryList = router.New(
		TestQueryList,
		router.Summary("Test query list"),
		router.Description("Test query list model"),
		router.Security(&security.Basic{}),
		router.Responses(router.Response{
			"200": router.ResponseItem{
				Model: []TestQueryListReq{},
			},
		}),
	)
	noModel = router.NewX(
		TestNoModel,
		router.Summary("Test no model"),
		router.Description("Test no model"),
		router.Responses(router.Response{
			"200": router.ResponseItem{
				Description: "success",
			},
		}),
	)
	queryPath = router.New(
		TestQueryPath,
		router.Summary("Test query path"),
		router.Description("Test query path model"),
		router.Responses(router.Response{
			"200": router.ResponseItem{
				Description: "success",
				Model:       TestQueryPathReq{},
			},
		}),
	)
	formEncode = router.New(
		TestForm,
		router.Summary("Test form"),
		router.ContentType(fiber.MIMEApplicationForm, router.ContentTypeRequest),
	)
	body = router.New(
		TestJson,
		router.Summary("Test json body"),
		router.Responses(router.Response{
			"200": router.ResponseItem{
				Model: TestFormReq{},
			},
		}),
	)
	file = router.New(
		TestFile,
		router.Summary("Test file upload"),
		router.ContentType(fiber.MIMEApplicationForm, router.ContentTypeRequest),
	)
)
