package openapi

import (
	"api/http_server/http_util"
	"api/http_server/middleware"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewOpenApi3Router(config Config) (func(router chi.Router), error) {
	basicAuth := middleware.BasicAuth(config.BasicAuthUsername, config.BasicAuthPassword, config.BasicAuthRealm)

	swaggerJson, err := newOpenApi3(OpenApi3Config{
		DomainWithProtocol:         config.Domain,
		OAuth2TokenUrl:             config.OAuth2TokenUrl,
		OAuth2AuthorizationCodeUrl: config.OAuth2AuthorizationCodeUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating openapi3: %w", err)
	}

	return func(r chi.Router) {
		FileServer(r, "/", http.FS(openapi3content))

		r.Get("/swagger.json", basicAuth(func(w http.ResponseWriter, req *http.Request) {
			http_util.WriteJson(w, http.StatusOK, swaggerJson)
		}))

		r.Get("/swagger-config.json", basicAuth(func(w http.ResponseWriter, req *http.Request) {
			http_util.WriteJson(w, http.StatusOK, SwaggerUIConfig{
				Url:               config.Domain + "/docs/swagger.json",
				OAuth2RedirectUrl: config.Domain + "/docs/oauth2-redirect.html",
				DomId:             "#swagger-ui",
				DeepLinking:       true,
				ValidatorUrl:      nil,
			})
		}))
	}, nil
}
