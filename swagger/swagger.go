package swagger

import (
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/long2ice/fibers/constants"
	"github.com/long2ice/fibers/router"
	"github.com/long2ice/fibers/security"
	log "github.com/sirupsen/logrus"
)

type Swagger struct {
	Title          string
	Description    string
	Version        string
	DocsUrl        string
	RedocUrl       string
	OpenAPIUrl     string
	Routers        map[string]map[string]*router.Router
	Servers        openapi3.Servers
	TermsOfService string
	Contact        *openapi3.Contact
	License        *openapi3.License
	OpenAPI        *openapi3.T
	SwaggerOptions map[string]interface{}
	RedocOptions   map[string]interface{}
}

func New(title, description, version string, options ...Option) *Swagger {
	swagger := &Swagger{
		Title:       title,
		Description: description,
		Version:     version,
		DocsUrl:     "/docs",
		RedocUrl:    "/redoc",
		OpenAPIUrl:  "/openapi.json",
	}
	for _, option := range options {
		option(swagger)
	}
	return swagger
}

func (swagger *Swagger) getSecurityRequirements(
	securities []security.ISecurity,
) *openapi3.SecurityRequirements {
	securityRequirements := openapi3.NewSecurityRequirements()
	for _, s := range securities {
		provide := string(s.Provider())
		swagger.OpenAPI.Components.SecuritySchemes[provide] = &openapi3.SecuritySchemeRef{
			Value: s.Scheme(),
		}
		securityRequirements.With(openapi3.NewSecurityRequirement().Authenticate(provide))
	}
	return securityRequirements
}

func (swagger *Swagger) getSchemaByType(t interface{}, request bool) *openapi3.Schema {
	var schema *openapi3.Schema
	var m float64
	m = float64(0)
	switch t.(type) {
	case int, int8, int16, *int, *int8, *int16:
		schema = openapi3.NewIntegerSchema()
	case uint, uint8, uint16, *uint, *uint8, *uint16:
		schema = openapi3.NewIntegerSchema()
		schema.Min = &m
	case int32, *int32:
		schema = openapi3.NewInt32Schema()
	case uint32, *uint32:
		schema = openapi3.NewInt32Schema()
		schema.Min = &m
	case int64, *int64:
		schema = openapi3.NewInt64Schema()
	case uint64, *uint64:
		schema = openapi3.NewInt64Schema()
		schema.Min = &m
	case string, *string:
		schema = openapi3.NewStringSchema()
	case time.Time, *time.Time:
		schema = openapi3.NewDateTimeSchema()
	case uuid.UUID, *uuid.UUID:
		schema = openapi3.NewUUIDSchema()
	case float32, float64, *float32, *float64:
		schema = openapi3.NewFloat64Schema()
	case bool, *bool:
		schema = openapi3.NewBoolSchema()
	case []byte:
		schema = openapi3.NewBytesSchema()
	case *multipart.FileHeader:
		schema = openapi3.NewStringSchema()
		schema.Format = "binary"
	case []*multipart.FileHeader:
		schema = openapi3.NewArraySchema()
		schema.Items = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:   "string",
				Format: "binary",
			},
		}
	default:
		if request {
			schema = swagger.getRequestSchemaByModel(t)
		} else {
			schema = swagger.getResponseSchemaByModel(t)
		}
	}
	return schema
}

func (swagger *Swagger) getRequestSchemaByModel(model interface{}) *openapi3.Schema {
	type_ := reflect.TypeOf(model)
	value_ := reflect.ValueOf(model)
	schema := openapi3.NewObjectSchema()
	if type_.Kind() == reflect.Ptr {
		type_ = type_.Elem()
	}
	if value_.Kind() == reflect.Ptr {
		value_ = value_.Elem()
	}
	if type_.Kind() == reflect.Struct {
		for i := 0; i < type_.NumField(); i++ {
			field := type_.Field(i)
			value := value_.Field(i)
			tags, err := structtag.Parse(string(field.Tag))
			if err != nil {
				log.Fatal(err)
			}
			_, err = tags.Get(constants.EMBED)
			if err == nil {
				embedSchema := swagger.getRequestSchemaByModel(value.Interface())
				for key, embedProperty := range embedSchema.Properties {
					schema.Properties[key] = embedProperty
				}
				for _, name := range embedSchema.Required {
					schema.Required = append(schema.Required, name)
				}
			}
			tag, err := tags.Get(constants.FORM)
			if err != nil {
				tag, err = tags.Get(constants.JSON)
				if err != nil {
					continue
				}
			}
			fieldSchema := swagger.getSchemaByType(value.Interface(), true)
			schema.Properties[tag.Name] = openapi3.NewSchemaRef("", fieldSchema)
			validateTag, err := tags.Get(constants.VALIDATE)
			if err == nil {
				if validateTag.Name == "required" {
					schema.Required = append(schema.Required, tag.Name)
				}
				options := validateTag.Options
				if len(options) > 0 {
					schema.Properties[tag.Name] = swagger.getValidateSchemaByOptions(
						value.Interface(),
						options,
					)
					fieldSchema = schema.Properties[tag.Name].Value
				}
			}
			defaultTag, err := tags.Get(constants.DEFAULT)
			if err == nil {
				fieldSchema.Default = defaultTag.Name
			}
			exampleTag, err := tags.Get(constants.EXAMPLE)
			if err == nil {
				fieldSchema.Example = exampleTag.Name
			}
			descriptionTag, err := tags.Get(constants.DESCRIPTION)
			if err == nil {
				fieldSchema.Description = descriptionTag.Name
			}
		}
	} else if type_.Kind() == reflect.Slice {
		schema = openapi3.NewArraySchema()
		schema.Items = &openapi3.SchemaRef{Value: swagger.getRequestSchemaByModel(reflect.New(type_.Elem()).Elem().Interface())}
	} else if type_.Kind() != reflect.Map {
		schema = swagger.getSchemaByType(model, true)
	}
	return schema
}

func (swagger *Swagger) getRequestBodyByModel(
	model interface{},
	contentType string,
) *openapi3.RequestBodyRef {
	body := &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody(),
	}
	if model == nil {
		return body
	}
	schema := swagger.getRequestSchemaByModel(model)
	body.Value.Required = true
	if contentType == "" {
		contentType = fiber.MIMEApplicationJSON
	}
	body.Value.Content = openapi3.NewContentWithSchema(schema, []string{contentType})
	return body
}

func (swagger *Swagger) getResponseSchemaByModel(model interface{}) *openapi3.Schema {
	schema := openapi3.NewObjectSchema()
	if model == nil {
		return schema
	}
	type_ := reflect.TypeOf(model)
	value_ := reflect.ValueOf(model)
	if type_.Kind() == reflect.Ptr {
		type_ = type_.Elem()
	}
	if value_.Kind() == reflect.Ptr {
		value_ = value_.Elem()
	}
	if type_.Kind() == reflect.Struct {
		for i := 0; i < type_.NumField(); i++ {
			field := type_.Field(i)
			value := value_.Field(i)
			fieldSchema := swagger.getSchemaByType(value.Interface(), false)
			tags, err := structtag.Parse(string(field.Tag))
			if err != nil {
				panic(err)
			}
			_, err = tags.Get(constants.EMBED)
			if err == nil {
				embedSchema := swagger.getResponseSchemaByModel(value.Interface())
				for key, embedProperty := range embedSchema.Properties {
					schema.Properties[key] = embedProperty
				}
				for _, name := range embedSchema.Required {
					schema.Required = append(schema.Required, name)
				}
			}
			tag, err := tags.Get(constants.JSON)
			if err != nil {
				continue
			}
			validateTag, err := tags.Get(constants.VALIDATE)
			if err == nil && validateTag.Name == "required" {
				schema.Required = append(schema.Required, tag.Name)
			}
			descriptionTag, err := tags.Get(constants.DESCRIPTION)
			if err == nil {
				fieldSchema.Description = descriptionTag.Name
			}
			defaultTag, err := tags.Get(constants.DEFAULT)
			if err == nil {
				fieldSchema.Default = defaultTag.Name
			}
			exampleTag, err := tags.Get(constants.EXAMPLE)
			if err == nil {
				fieldSchema.Example = exampleTag.Name
			}
			schema.Properties[tag.Name] = openapi3.NewSchemaRef("", fieldSchema)
		}
	} else if type_.Kind() == reflect.Slice {
		schema = openapi3.NewArraySchema()
		schema.Items = &openapi3.SchemaRef{Value: swagger.getResponseSchemaByModel(reflect.New(type_.Elem()).Elem().Interface())}
	} else if type_.Kind() != reflect.Map {
		schema = swagger.getSchemaByType(model, false)
	}
	return schema
}

func (swagger *Swagger) getResponses(
	response router.Response,
	contentType string,
) openapi3.Responses {
	ret := openapi3.NewResponses()
	for k, v := range response {
		schema := swagger.getResponseSchemaByModel(v.Model)
		var content openapi3.Content
		if contentType == "" || contentType == fiber.MIMEApplicationJSON {
			content = openapi3.NewContentWithJSONSchema(schema)
		} else {
			content = openapi3.NewContentWithSchema(schema, []string{contentType})
		}
		description := v.Description
		ret[k] = &openapi3.ResponseRef{
			Value: &openapi3.Response{
				Description: &description,
				Content:     content,
				Headers:     v.Headers,
			},
		}
	}
	return ret
}

func (swagger *Swagger) getValidateSchemaByOptions(
	value interface{},
	options []string,
) *openapi3.SchemaRef {
	schema := openapi3.NewSchemaRef("", swagger.getSchemaByType(value, true))
	for _, option := range options {
		if strings.HasPrefix(option, "oneof=") {
			optionItems := strings.Split(option[6:], " ")
			enums := make([]interface{}, len(optionItems))
			for i, optionItem := range optionItems {
				enums[i] = optionItem
			}
			schema.Value.WithEnum(enums...)
		}
		if strings.HasPrefix(option, "max=") {
			value, err := strconv.ParseFloat(option[4:], 64)
			if err != nil {
				log.Panicln(err)
			}
			schema.Value.WithMax(value)
		}
		if strings.HasPrefix(option, "min=") {
			value, err := strconv.ParseFloat(option[4:], 64)
			if err != nil {
				log.Panicln(err)
			}
			schema.Value.WithMin(value)
		}
		if strings.HasPrefix(option, "len=") {
			value, err := strconv.ParseInt(option[4:], 10, 64)
			if err != nil {
				log.Panicln(err)
			}
			schema.Value.WithLength(value)
		}
	}
	return schema
}

func (swagger *Swagger) getParametersByModel(model interface{}) openapi3.Parameters {
	parameters := openapi3.NewParameters()
	if model == nil {
		return parameters
	}
	type_ := reflect.TypeOf(model)
	if type_.Kind() == reflect.Ptr {
		type_ = type_.Elem()
	}
	value_ := reflect.ValueOf(model)
	if value_.Kind() == reflect.Ptr {
		value_ = value_.Elem()
	}
	for i := 0; i < type_.NumField(); i++ {
		field := type_.Field(i)
		value := value_.Field(i)
		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			panic(err)
		}
		_, err = tags.Get(constants.EMBED)
		if err == nil {
			embedParameters := swagger.getParametersByModel(value.Interface())
			for _, embedParameter := range embedParameters {
				parameters = append(parameters, embedParameter)
			}
		}
		parameter := &openapi3.Parameter{
			Schema: openapi3.NewSchemaRef("", swagger.getSchemaByType(value.Interface(), true)),
		}
		queryTag, err := tags.Get(constants.QUERY)
		if err == nil {
			parameter.In = openapi3.ParameterInQuery
			parameter.Name = queryTag.Name
		}
		uriTag, err := tags.Get(constants.URI)
		if err == nil {
			parameter.In = openapi3.ParameterInPath
			parameter.Name = uriTag.Name
		}
		headerTag, err := tags.Get(constants.HEADER)
		if err == nil {
			parameter.In = openapi3.ParameterInHeader
			parameter.Name = headerTag.Name
		}
		cookieTag, err := tags.Get(constants.COOKIE)
		if err == nil {
			parameter.In = openapi3.ParameterInCookie
			parameter.Name = cookieTag.Name
		}
		if parameter.In == "" {
			continue
		}
		descriptionTag, err := tags.Get(constants.DESCRIPTION)
		if err == nil {
			parameter.WithDescription(descriptionTag.Name)
		}
		validateTag, err := tags.Get(constants.VALIDATE)
		if err == nil {
			parameter.WithRequired(validateTag.Name == "required")
			options := validateTag.Options
			if len(options) > 0 {
				parameter.Schema = swagger.getValidateSchemaByOptions(value.Interface(), options)
			}
		}
		defaultTag, err := tags.Get(constants.DEFAULT)
		if err == nil {
			parameter.Schema.Value.WithDefault(defaultTag.Name)
		}
		exampleTag, err := tags.Get(constants.EXAMPLE)
		if err == nil {
			parameter.Schema.Value.Example = exampleTag.Name
		}
		parameters = append(parameters, &openapi3.ParameterRef{
			Value: parameter,
		})
	}
	return parameters
}

// /:id -> /{id}
func (swagger *Swagger) fixPath(path string) string {
	reg := regexp.MustCompile("/:(\\w+)")
	return reg.ReplaceAllString(path, "/{${1}}")
}

func (swagger *Swagger) getPaths() openapi3.Paths {
	paths := make(openapi3.Paths)
	for path, m := range swagger.Routers {
		pathItem := &openapi3.PathItem{}
		for method, r := range m {
			if r.Exclude {
				continue
			}
			model := r.Model
			operation := &openapi3.Operation{
				Tags:        r.Tags,
				OperationID: r.OperationID,
				Summary:     r.Summary,
				Description: r.Description,
				Deprecated:  r.Deprecated,
				Responses:   swagger.getResponses(r.Response, r.ResponseContentType),
				Parameters:  swagger.getParametersByModel(model),
				Security:    swagger.getSecurityRequirements(r.Securities),
			}
			requestBody := swagger.getRequestBodyByModel(model, r.RequestContentType)
			if method == http.MethodGet {
				pathItem.Get = operation
			} else if method == http.MethodPost {
				pathItem.Post = operation
				operation.RequestBody = requestBody
			} else if method == http.MethodDelete {
				pathItem.Delete = operation
			} else if method == http.MethodPut {
				pathItem.Put = operation
				operation.RequestBody = requestBody
			} else if method == http.MethodPatch {
				pathItem.Patch = operation
			} else if method == http.MethodHead {
				pathItem.Head = operation
			} else if method == http.MethodOptions {
				pathItem.Options = operation
			} else if method == http.MethodConnect {
				pathItem.Connect = operation
			} else if method == http.MethodTrace {
				pathItem.Trace = operation
			}
		}
		paths[swagger.fixPath(path)] = pathItem
	}
	return paths
}

func (swagger *Swagger) BuildOpenAPI() {
	components := openapi3.NewComponents()
	components.SecuritySchemes = openapi3.SecuritySchemes{}
	swagger.OpenAPI = &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          swagger.Title,
			Description:    swagger.Description,
			TermsOfService: swagger.TermsOfService,
			Contact:        swagger.Contact,
			License:        swagger.License,
			Version:        swagger.Version,
		},
		Servers:    swagger.Servers,
		Components: &components,
	}
	swagger.OpenAPI.Paths = swagger.getPaths()
}

func (swagger *Swagger) MarshalJSON() ([]byte, error) {
	return swagger.OpenAPI.MarshalJSON()
}

func (swagger *Swagger) WithDocsUrl(url string) *Swagger {
	DocsUrl(url)(swagger)
	return swagger
}

func (swagger *Swagger) WithRedocUrl(url string) *Swagger {
	RedocUrl(url)(swagger)
	return swagger
}

func (swagger *Swagger) WithTitle(title string) *Swagger {
	Title(title)(swagger)
	return swagger
}

func (swagger *Swagger) WithDescription(description string) *Swagger {
	Description(description)(swagger)
	return swagger
}

func (swagger *Swagger) WithVersion(version string) *Swagger {
	Version(version)(swagger)
	return swagger
}

func (swagger *Swagger) WithOpenAPIUrl(url string) *Swagger {
	OpenAPIUrl(url)(swagger)
	return swagger
}

func (swagger *Swagger) WithTermsOfService(termsOfService string) *Swagger {
	TermsOfService(termsOfService)(swagger)
	return swagger
}

func (swagger *Swagger) WithContact(contact *openapi3.Contact) *Swagger {
	Contact(contact)(swagger)
	return swagger
}

func (swagger *Swagger) WithLicense(license *openapi3.License) *Swagger {
	License(license)(swagger)
	return swagger
}

func (swagger *Swagger) WithServers(servers []*openapi3.Server) *Swagger {
	Servers(servers)(swagger)
	return swagger
}

func (swagger *Swagger) WithSwaggerOptions(options map[string]interface{}) *Swagger {
	SwaggerOptions(options)(swagger)
	return swagger
}

func (swagger *Swagger) WithRedocOptions(options map[string]interface{}) *Swagger {
	RedocOptions(options)(swagger)
	return swagger
}
