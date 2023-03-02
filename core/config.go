package core

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	AwsUserPoolId               string
	AwsRegion                   string
	DatabaseUrl                 string
	ImagesApiDomain             string
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

	c.AwsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	if c.AwsAccessKeyId == "" {
		return errors.New("missing env AWS_ACCESS_KEY_ID")
	}

	c.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	if c.AwsSecretAccessKey == "" {
		return errors.New("missing env AWS_SECRET_ACCESS_KEY")
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
