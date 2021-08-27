package http_transport

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"simple_gopher"
)

func openApi3Router(config simple_gopher.Config) (func(router chi.Router), error) {
	basicAuth := newBasicAuth(config.BasicAuthUsername, config.BasicAuthPassword, config.BasicAuthRealm)

	swaggerJson, err := newOpenApi3(OpenApi3Config{
		DomainWithProtocol:         config.Domain,
		OAuth2TokenUrl:             config.OAuth2TokenUrl,
		OAuth2AuthorizationCodeUrl: config.OAuth2AuthorizationCodeUrl,
	})
	if err != nil {
		return nil, err
	}

	return func(r chi.Router) {
		FileServer(r, "/", http.FS(openapi3content))

		r.Get("/swagger.json", basicAuth(func(w http.ResponseWriter, req *http.Request) {
			respondJson(w, http.StatusOK, swaggerJson)
		}))

		r.Get("/swagger-config.json", basicAuth(func(w http.ResponseWriter, req *http.Request) {
			respondJson(w, http.StatusOK, SwaggerUIConfig{
				Url:               config.Domain + "/docs/swagger.json",
				OAuth2RedirectUrl: config.Domain + "/docs/oauth2-redirect.html",
				DomId:             "#swagger-ui",
				DeepLinking:       true,
				ValidatorUrl:      nil,
			})
		}))
	}, nil
}
