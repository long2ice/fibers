package swagger

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Option func(swagger *Swagger)

func DocsUrl(url string) Option {
	return func(swagger *Swagger) {
		swagger.DocsUrl = url
	}
}
func RedocUrl(url string) Option {
	return func(swagger *Swagger) {
		swagger.RedocUrl = url
	}
}
func Title(title string) Option {
	return func(swagger *Swagger) {
		swagger.Title = title
	}
}
func Description(description string) Option {
	return func(swagger *Swagger) {
		swagger.Description = description
	}
}
func Version(version string) Option {
	return func(swagger *Swagger) {
		swagger.Version = version
	}
}
func OpenAPIUrl(url string) Option {
	return func(swagger *Swagger) {
		swagger.OpenAPIUrl = url
	}
}
func Servers(servers openapi3.Servers) Option {
	return func(swagger *Swagger) {
		swagger.Servers = servers
	}
}
func TermsOfService(TermsOfService string) Option {
	return func(swagger *Swagger) {
		swagger.TermsOfService = TermsOfService
	}
}
func Contact(Contact *openapi3.Contact) Option {
	return func(swagger *Swagger) {
		swagger.Contact = Contact
	}
}
func License(License *openapi3.License) Option {
	return func(swagger *Swagger) {
		swagger.License = License
	}
}
func SwaggerOptions(options map[string]interface{}) Option {
	return func(swagger *Swagger) {
		swagger.SwaggerOptions = options
	}
}
func RedocOptions(options map[string]interface{}) Option {
	return func(swagger *Swagger) {
		swagger.RedocOptions = options
	}
}
