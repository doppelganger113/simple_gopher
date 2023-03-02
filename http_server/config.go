package http_server

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                       uint
	CorsAllowOrigins           []string
	Timeout                    time.Duration
	ReadTimeout                time.Duration
	WriteTimeout               time.Duration
	IdleTimeout                time.Duration
	ReadHeaderTimeout          time.Duration
	HeartbeatUrl               string
	DebugRoutes                bool
	BasicAuthUsername          string
	BasicAuthPassword          string
	BasicAuthRealm             string
	Domain                     string
	OAuth2TokenUrl             string
	OAuth2AuthorizationCodeUrl string
}

func NewDefaultConfig() Config {
	return Config{
		Port:              3000,
		Timeout:           30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		HeartbeatUrl:      "/",
		BasicAuthUsername: "admin",
		BasicAuthPassword: "1234",
		BasicAuthRealm:    "simple_gopher",
	}
}

func (c *Config) LoadFromEnv() error {
	if domain := os.Getenv("DOMAIN"); domain != "" {
		c.Domain = domain
	}

	if oauthTokenUrl := os.Getenv("OAUTH2_TOKEN_URL"); oauthTokenUrl != "" {
		c.OAuth2TokenUrl = oauthTokenUrl
	}

	if code := os.Getenv("OAUTH_AUTHORIZATION_CODE_URL"); code != "" {
		c.OAuth2AuthorizationCodeUrl = code
	}

	if realm := os.Getenv("BASIC_AUTH_REALM"); realm != "" {
		c.BasicAuthRealm = realm
	}

	if username := os.Getenv("BASIC_AUTH_USERNAME"); username != "" {
		c.BasicAuthUsername = username
	}

	if pass := os.Getenv("BASIC_AUTH_PASSWORD"); pass != "" {
		c.BasicAuthPassword = pass
	}

	if corsOrigins := os.Getenv("CORS_ALLOW_ORIGINS"); corsOrigins != "" {
		c.CorsAllowOrigins = extractCorsOrigins(corsOrigins)
	}

	// Optional
	if port := os.Getenv("PORT"); port != "" {
		parsed, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			return err
		}
		c.Port = uint(parsed)
	}

	if debugRoutes := os.Getenv("DEBUG_ROUTES"); debugRoutes == "true" {
		c.DebugRoutes = true
	}

	if heartbeatUrl := os.Getenv("HEARTBEAT_URL"); heartbeatUrl != "" {
		c.HeartbeatUrl = heartbeatUrl
	}

	if timeout := os.Getenv("TIMEOUT"); timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}
		c.Timeout = parsed
	}

	if timeout := os.Getenv("TIMEOUT"); timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err != nil {
			return err
		}
		c.Timeout = parsed
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Domain == "" {
		return errors.New("missing env DOMAIN")
	}

	if c.OAuth2TokenUrl == "" {
		return errors.New("missing env OAUTH2_TOKEN_URL")
	}

	if c.OAuth2AuthorizationCodeUrl == "" {
		return errors.New("missing env OAUTH_AUTHORIZATION_CODE_URL")
	}

	if c.BasicAuthRealm == "" {
		c.BasicAuthRealm = "Forbidden"
	}

	if c.BasicAuthUsername == "" {
		return errors.New("missing env BASIC_AUTH_USERNAME")
	}

	if c.BasicAuthPassword == "" {
		return errors.New("missing env BASIC_AUTH_PASSWORD")
	}

	if len(c.CorsAllowOrigins) == 0 {
		return errors.New("empty CORS_ALLOW_ORIGINS")
	}

	if c.Port < 80 {
		return errors.New("invalid PORT value")
	}

	if c.HeartbeatUrl[0] != '/' {
		return errors.New("env HEARTBEAT_URL must start with /")
	}

	return nil
}

func extractCorsOrigins(corsOrigins string) []string {
	origins := strings.Split(corsOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return origins
}
