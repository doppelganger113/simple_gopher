package image_resize_api

import (
	"fmt"
	"net/http"
	"time"
)

type ResizeApi struct {
	domain string
	client *http.Client
}

func NewResizeApi(config Config) *ResizeApi {
	return &ResizeApi{
		domain: config.ImagesApiDomain,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (api *ResizeApi) url(relativePath string) string {
	return fmt.Sprintf("%s%s", api.domain, relativePath)
}

func isResponseOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
