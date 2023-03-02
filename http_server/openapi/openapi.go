package openapi

import (
	"api/http_server/http_util"
	"embed"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

//go:embed docs
var openapi3content embed.FS

type OpenApi3Config struct {
	DomainWithProtocol         string
	OAuth2TokenUrl             string
	OAuth2AuthorizationCodeUrl string
}

type SwaggerUIConfig struct {
	// Swagger json url
	Url               string  `json:"url"`
	OAuth2RedirectUrl string  `json:"oauth2RedirectUrl"`
	DomId             string  `json:"dom_id"`
	DeepLinking       bool    `json:"deepLinking"`
	ValidatorUrl      *string `json:"validator_url"`
}

func newOpenApi3(config OpenApi3Config) (*openapi3.T, error) {
	swagger := &openapi3.T{OpenAPI: "3.0.0"}
	swagger.Info = &openapi3.Info{
		Title:          "Golang API",
		Description:    "Prototype API. When performing authentication use only the client_id, you won't need the secret",
		TermsOfService: "",
		License: &openapi3.License{
			Name: "This API is prohibited for external use.",
		},
		Version: "0.0.0",
	}

	swagger.Servers = openapi3.Servers{
		&openapi3.Server{
			URL:         config.DomainWithProtocol,
			Description: "Golang API",
		},
	}

	errResponseSchemaRef, _, err := openapi3gen.NewSchemaRefForValue(&http_util.FailureResponse{})
	if err != nil {
		return nil, err
	}

	swagger.Components.Schemas = openapi3.Schemas{
		"Image": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Properties: map[string]*openapi3.SchemaRef{
					"name": {
						Value: &openapi3.Schema{Type: "string", Example: "my plane"},
					},
					"format": {
						Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"jpg", "png", "webp"}},
					},
					"originalFile": {
						Value: &openapi3.Schema{Type: "string", Format: "binary"},
					},
					"croppedFile": {
						Value: &openapi3.Schema{Type: "string", Format: "binary"},
					},
				},
				Required: []string{"name", "format", "originalFile", "croppedFile"},
			},
		},
		"ErrResponse": errResponseSchemaRef,
		"CreateImage": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: "object",
				Properties: map[string]*openapi3.SchemaRef{
					"name": {
						Value: &openapi3.Schema{Type: "string", Example: "my plane"},
					},
					"format": {
						Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"jpg", "png", "webp"}},
					},
					"originalFile": {
						Value: &openapi3.Schema{Type: "string", Format: "binary"},
					},
					"croppedFile": {
						Value: &openapi3.Schema{Type: "string", Format: "binary"},
					},
				},
				Required: []string{"name", "format", "originalFile", "croppedFile"},
			},
		},
	}

	swagger.Components.RequestBodies = openapi3.RequestBodies{
		"UploadNewImage": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Create a new image. Ensure that the cropped image is in one of the allowed aspect ratios: `1:1` `3:2` `4:3` `5:8` `16:9`, otherwise it will fail.").
				WithRequired(true).
				WithContent(openapi3.NewContentWithFormDataSchemaRef(
					&openapi3.SchemaRef{
						Ref: "#/components/schemas/CreateImage",
					},
				)),
		},
		"UpdateImage": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription(
					"If you upload images, you will need to provide cropped, original and format. Name can be standalone",
				).
				WithRequired(true).
				WithContent(openapi3.NewContentWithFormDataSchemaRef(&openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "object",
						Properties: map[string]*openapi3.SchemaRef{
							"name": {
								Value: &openapi3.Schema{Type: "string", Example: "my plane"},
							},
							"format": {
								Value: &openapi3.Schema{Type: "string", Enum: []interface{}{"jpg", "png", "webp"}},
							},
							"originalFile": {
								Value: &openapi3.Schema{Type: "string", Format: "binary"},
							},
							"croppedFile": {
								Value: &openapi3.Schema{Type: "string", Format: "binary"},
							},
						},
					},
				})),
		},
	}

	swagger.Components.Responses = openapi3.Responses{
		"ImageResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Created image where sizes in size map may be nullable").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Ref: "#/components/schemas/Image",
						},
					),
				),
		},
		"ImagesResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Images").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "array",
								Items: &openapi3.SchemaRef{
									Ref: "#/components/schemas/Image",
								},
							},
						},
					),
				),
		},
		"EmptyResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Ok empty response"),
		},
		"BadRequestResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Bad request").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Ref: "#/components/schemas/ErrResponse",
						},
					),
				),
		},
		"ForbiddenResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Forbidden").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Ref: "#/components/schemas/ErrResponse",
						},
					),
				),
		},
		"UnauthorizedResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Unauthorized").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Ref: "#/components/schemas/ErrResponse",
						},
					),
				),
		},
		"NotFoundResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Resource not found").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Ref: "#/components/schemas/ErrResponse",
						},
					),
				),
		},
		"ServerErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("ServerErrorResponse").
				WithContent(
					openapi3.NewContentWithJSONSchemaRef(
						&openapi3.SchemaRef{
							Ref: "#/components/schemas/ErrResponse",
						},
					),
				),
		},
	}

	swagger.Paths = openapi3.Paths{
		"/api/v1/images": &openapi3.PathItem{
			Summary: "Images aka pictures",
			Get: &openapi3.Operation{
				OperationID: "GetImages",
				Tags:        []string{"Images"},
				Description: "Fetch list of images",
				Parameters: openapi3.Parameters{
					{
						Value: &openapi3.Parameter{
							Name: "size",
							In:   "query",
							Description: fmt.Sprintf(
								"Number of results, default is %d and maximum is %d",
								20, 100,
							),
						},
					},
					{
						Value: &openapi3.Parameter{
							Name:        "page",
							In:          "query",
							Description: "Page number for pagination, minimum 1",
						},
					},
					{
						Value: &openapi3.Parameter{
							Name:        "order",
							In:          "query",
							Description: "Specify descending or ascending order",
							Schema: &openapi3.SchemaRef{
								Value: openapi3.NewSchema().WithEnum("First", "Second"),
							},
						},
					},
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/ImagesResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ServerErrorResponse",
					},
				},
			},
			Post: &openapi3.Operation{
				OperationID: "Upload and create an image",
				Tags:        []string{"Images"},
				Description: "Upload and save the image, requires admin authorization",
				Security: &openapi3.SecurityRequirements{
					openapi3.SecurityRequirement{
						"oauth2": []string{},
					},
				},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/UploadNewImage",
				},
				Responses: openapi3.Responses{
					"201": &openapi3.ResponseRef{
						Ref: "#/components/responses/ImageResponse",
					},
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/BadRequestResponse",
					},
					"401": &openapi3.ResponseRef{
						Ref: "#/components/responses/ForbiddenResponse",
					},
					"403": &openapi3.ResponseRef{
						Ref: "#/components/responses/UnauthorizedResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ServerErrorResponse",
					},
				},
			},
		},
		"/api/v1/images/{id}": &openapi3.PathItem{
			Summary: "Image",
			Get: &openapi3.Operation{
				OperationID: "GetImage",
				Tags:        []string{"Images"},
				Description: "Fetch image info",
				Parameters: openapi3.Parameters{
					{
						Value: &openapi3.Parameter{
							Name:        "id",
							In:          "path",
							Description: "Id of image",
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type:   "string",
									Format: "uuid",
								},
							},
						},
					},
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/ImageResponse",
					},
					"404": &openapi3.ResponseRef{
						Ref: "#/components/responses/NotFoundResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ServerErrorResponse",
					},
				},
			},
			Delete: &openapi3.Operation{
				OperationID: "DeleteImage",
				Tags:        []string{"Images"},
				Description: "Delete image and invalidate CDN images",
				Security: &openapi3.SecurityRequirements{
					openapi3.SecurityRequirement{
						"oauth2": []string{},
					},
				},
				Parameters: openapi3.Parameters{
					{
						Value: &openapi3.Parameter{
							Name:        "id",
							In:          "path",
							Description: "Id of image",
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type:   "string",
									Format: "uuid",
								},
							},
						},
					},
				},
				Responses: openapi3.Responses{
					"204": &openapi3.ResponseRef{
						Ref: "#/components/responses/EmptyResponse",
					},
					"404": &openapi3.ResponseRef{
						Ref: "#/components/responses/NotFoundResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ServerErrorResponse",
					},
				},
			},
			Patch: &openapi3.Operation{
				OperationID: "Update image",
				Tags:        []string{"Images"},
				Description: "Update existing image or change the name. Note that this will invalidate the cashed image on edge locations.",
				Security: &openapi3.SecurityRequirements{
					openapi3.SecurityRequirement{
						"oauth2": []string{},
					},
				},
				Parameters: openapi3.Parameters{
					{
						Value: &openapi3.Parameter{
							Name:        "id",
							In:          "path",
							Description: "Id of image",
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type:   "string",
									Format: "uuid",
								},
							},
						},
					},
				},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/UpdateImage",
				},
				Responses: openapi3.Responses{
					"200": &openapi3.ResponseRef{
						Ref: "#/components/responses/ImageResponse",
					},
					"400": &openapi3.ResponseRef{
						Ref: "#/components/responses/BadRequestResponse",
					},
					"401": &openapi3.ResponseRef{
						Ref: "#/components/responses/ForbiddenResponse",
					},
					"403": &openapi3.ResponseRef{
						Ref: "#/components/responses/UnauthorizedResponse",
					},
					"404": &openapi3.ResponseRef{
						Ref: "#/components/responses/NotFoundResponse",
					},
					"500": &openapi3.ResponseRef{
						Ref: "#/components/responses/ServerErrorResponse",
					},
				},
			},
		},
	}

	swagger.Components.SecuritySchemes = openapi3.SecuritySchemes{
		"oauth2": &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type: "oauth2",
				Flows: &openapi3.OAuthFlows{
					AuthorizationCode: &openapi3.OAuthFlow{
						AuthorizationURL: config.OAuth2AuthorizationCodeUrl,
						TokenURL:         config.OAuth2TokenUrl,
						RefreshURL:       config.OAuth2TokenUrl,
					},
				},
			},
		},
	}

	return swagger, nil
}
