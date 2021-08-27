package cognito

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
	"simple_gopher/cloud_patterns"
	"simple_gopher/storage"
)

type AuthConsumer struct {
	config      Config
	sqsClient   *sqs.SQS
	closed      chan error
	cancel      context.CancelFunc
	userStorage storage.UserRepository
}

func NewCognitoAuthConsumer(
	userStorage storage.UserRepository, config Config,
) *AuthConsumer {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sqsClient := sqs.New(sess)

	return &AuthConsumer{
		config:      config,
		sqsClient:   sqsClient,
		userStorage: userStorage,
		closed:      make(chan error),
	}
}

func (authConsumer *AuthConsumer) StartConsuming(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := authConsumer.ConsumeMessages()
			if err != nil {
				log.Err(err).Caller().Msgf("error consuming messages")
			}

			err = cloud_patterns.SleepSecondsWithContext(ctx, authConsumer.config.SqsPostAuthIntervalSec)
			if err != nil {
				return err
			}
		}
	}
}

func (authConsumer *AuthConsumer) StartConsumingAsync(ctx context.Context) {
	log.Info().Msg("Started consuming authentications")

	derivedCtx, cancel := context.WithCancel(ctx)
	authConsumer.cancel = cancel
	go func() {
		authConsumer.closed <- authConsumer.StartConsuming(derivedCtx)
		close(authConsumer.closed)
	}()
}

func (authConsumer *AuthConsumer) Shutdown() error {
	log.Info().Msg("Shutting down")
	authConsumer.cancel()
	if err := <-authConsumer.closed; err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (authConsumer *AuthConsumer) deleteMessage(receiptHandle string) error {
	_, err := authConsumer.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(authConsumer.config.SqsPostAuthUrl),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func (authConsumer *AuthConsumer) ConsumeMessages() error {
	output, err := authConsumer.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames:        []*string{aws.String("All")},
		MaxNumberOfMessages:   aws.Int64(10),
		MessageAttributeNames: nil,
		QueueUrl:              aws.String(authConsumer.config.SqsPostAuthUrl),
		VisibilityTimeout:     aws.Int64(30),
		WaitTimeSeconds:       aws.Int64(20),
	})
	if err != nil {
		return err
	}

	for _, msg := range output.Messages {
		if err = authConsumer.handleMessage(msg); err != nil {
			log.Err(err).Caller().Msg("failed handling auth message")
		}
	}

	return nil
}

func (authConsumer *AuthConsumer) handleMessage(message *sqs.Message) error {
	postAuthEvents, parseErr := ParsePostAuthEvent(message)
	if parseErr != nil {
		return parseErr
	}

	log.Info().Caller().Msgf("Consuming %s", postAuthEvents.Request.UserAttributes.Email)

	newUser := storage.UserCreationDto{
		Email:       postAuthEvents.Request.UserAttributes.Email,
		Role:        storage.AuthRoleNone,
		CogUsername: postAuthEvents.Username,
		CogSub:      postAuthEvents.Request.UserAttributes.Sub,
		CogName:     postAuthEvents.Username,
		Disabled:    false,
	}
	_, err := authConsumer.userStorage.Create(context.Background(), newUser)
	if err != nil {
		if errors.Is(err, storage.DuplicateErr) {
			log.Info().Caller().Msgf("%s already exists, skipping", newUser.Email)
		} else {
			return err
		}
	}

	if err = authConsumer.deleteMessage(*message.ReceiptHandle); err != nil {
		return err
	}

	log.Info().Caller().Msgf("Successfully consumed %s", newUser.Email)

	return nil
}
