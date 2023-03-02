package resize

import (
	"api/core"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Client struct {
	domain string
	client *http.Client
	logger *zerolog.Logger
}

func NewClient(
	config core.Config,
	logger *zerolog.Logger,
) *Client {
	return &Client{
		domain: config.ImagesApiDomain,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (client *Client) url(relativePath string) string {
	return fmt.Sprintf("%s%s", client.domain, relativePath)
}

func isResponseOk(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
