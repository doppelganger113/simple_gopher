package cognito

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type PostAuthEvent struct {
	Version    string `json:"version"`
	Region     string `json:"region"`
	UserPoolId string `json:"userPoolId"`
	Username   string `json:"userName"`
	Request    struct {
		UserAttributes struct {
			Sub               string `json:"sub"`
			CognitoEmailAlias string `json:"cognito:email_alias"`
			CognitoUserStatus string `json:"cognito:user_status"`
			EmailVerified     string `json:"email_verified"`
			Email             string `json:"email"`
			Identities        string `json:"identities"`
		} `json:"userAttributes"`
	} `json:"request"`
}

func ParsePostAuthEvent(message *sqs.Message) (*PostAuthEvent, error) {
	var parsed PostAuthEvent
	if err := json.Unmarshal([]byte(*message.Body), &parsed); err != nil {
		return nil, err
	}
	return &parsed, nil
}
