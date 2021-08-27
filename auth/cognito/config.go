package cognito

type Config struct {
	SqsPostAuthIntervalSec uint
	SqsPostAuthUrl         string
	AwsRegion              string
	AwsUserPoolId          string
}
