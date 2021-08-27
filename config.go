package simple_gopher

import (
	"errors"
	"os"
	"simple_gopher/auth/cognito"
	"strconv"
	"strings"
)

type Config struct {
	Domain                      string
	OAuth2TokenUrl              string
	OAuth2AuthorizationCodeUrl  string
	BasicAuthUsername           string
	BasicAuthPassword           string
	BasicAuthRealm              string
	AwsUserPoolId               string
	AwsRegion                   string
	DatabaseUrl                 string
	ImagesApiDomain             string
	CorsAllowOrigins            []string
	PORT                        uint
	DebugRoutes                 bool
	AwsAccessKeyId              string
	AwsSecretAccessKey          string
	SqsPostAuthUrl              string
	SqsPostAuthIntervalSec      uint
	SqsPostAuthConsumerDisabled bool
}

func NewConfigFromEnv() (Config, error) {
	c := Config{}
	err := c.LoadFromEnvironment()
	return c, err
}

func (c *Config) LoadFromEnvironment() error {
	c.Domain = os.Getenv("DOMAIN")
	if c.Domain == "" {
		return errors.New("missing env DOMAIN")
	}

	c.OAuth2TokenUrl = os.Getenv("OAUTH2_TOKEN_URL")
	if c.OAuth2TokenUrl == "" {
		return errors.New("missing env OAUTH2_TOKEN_URL")
	}

	c.OAuth2AuthorizationCodeUrl = os.Getenv("OAUTH_AUTHORIZATION_CODE_URL")
	if c.OAuth2AuthorizationCodeUrl == "" {
		return errors.New("missing env OAUTH_AUTHORIZATION_CODE_URL")
	}

	c.BasicAuthRealm = os.Getenv("BASIC_AUTH_REALM")
	if c.BasicAuthRealm == "" {
		c.BasicAuthRealm = "Forbidden"
	}

	c.BasicAuthUsername = os.Getenv("BASIC_AUTH_USERNAME")
	if c.BasicAuthUsername == "" {
		return errors.New("missing env BASIC_AUTH_USERNAME")
	}

	c.BasicAuthPassword = os.Getenv("BASIC_AUTH_PASSWORD")
	if c.BasicAuthPassword == "" {
		return errors.New("missing env BASIC_AUTH_PASSWORD")
	}

	c.AwsRegion = os.Getenv("AWS_REGION")
	if c.AwsRegion == "" {
		return errors.New("missing env AWS_REGION")
	}

	c.AwsUserPoolId = os.Getenv("AWS_USER_POOL_ID")
	if c.AwsUserPoolId == "" {
		return errors.New("missing env AWS_USER_POOL_ID")
	}

	c.DatabaseUrl = os.Getenv("DATABASE_URL")
	if c.DatabaseUrl == "" {
		return errors.New("missing env DATABASE_URL")
	}

	c.ImagesApiDomain = os.Getenv("IMAGES_API_DOMAIN")
	if c.ImagesApiDomain == "" {
		return errors.New("missing env IMAGES_API_DOMAIN")
	}

	corsOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsOrigins == "" {
		return errors.New("missing env CORS_ALLOW_ORIGINS")
	}
	c.CorsAllowOrigins = extractCorsOrigins(corsOrigins)

	c.AwsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	if c.AwsAccessKeyId == "" {
		return errors.New("missing env AWS_ACCESS_KEY_ID")
	}

	c.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if c.AwsSecretAccessKey == "" {
		return errors.New("missing env AWS_SECRET_ACCESS_KEY")
	}

	// Optional
	if port := os.Getenv("PORT"); port != "" {
		parsed, err := strconv.ParseUint(port, 10, 64)
		if err != nil {
			return err
		}
		c.PORT = uint(parsed)
	} else {
		c.PORT = 3000
	}

	if debugRoutes := os.Getenv("DEBUG_ROUTES"); debugRoutes == "true" {
		c.DebugRoutes = true
	}

	c.SqsPostAuthUrl = os.Getenv("SQS_POST_AUTH_URL")
	if c.SqsPostAuthUrl == "" {
		return errors.New("missing env SQS_POST_AUTH_URL")
	}

	if seconds := os.Getenv("SQS_POST_AUTH_INTERVAL_SEC"); seconds != "" {
		parsedSeconds, err := strconv.Atoi(seconds)
		if err != nil {
			return err
		}
		c.SqsPostAuthIntervalSec = uint(parsedSeconds)
		if c.SqsPostAuthIntervalSec == 0 {
			return errors.New("env SQS_POST_AUTH_INTERVAL_SEC must not be 0")
		}
	} else {
		c.SqsPostAuthIntervalSec = 600
	}

	if os.Getenv("SQS_POST_AUTH_CONSUMER_DISABLED") == "true" {
		c.SqsPostAuthConsumerDisabled = true
	}

	return nil
}

func NewAuthConfig(config Config) cognito.Config {
	return cognito.Config{
		SqsPostAuthIntervalSec: config.SqsPostAuthIntervalSec,
		SqsPostAuthUrl:         config.SqsPostAuthUrl,
		AwsRegion:              config.AwsRegion,
		AwsUserPoolId:          config.AwsUserPoolId,
	}
}

func extractCorsOrigins(corsOrigins string) []string {
	origins := strings.Split(corsOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return origins
}
