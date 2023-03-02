package test

import (
	"api"
	"api/logger"
	"api/pkg/signaling"
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"net/http"
	"testing"
)

//func closeAfterTime(err chan<- error) {
//	time.Sleep(10 * time.Second)
//	err <- fmt.Errorf("closed")
//}

func getImages() (string, error) {
	resp, err := http.Get("http://localhost:3000/api/v1/images")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func TestIntegrationAppGetImages(t *testing.T) {
	SkipIfNotIntegrationTesting(t)

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5434/tcp", "5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "example",
			"POSTGRES_DB":       "db",
		},
	}
	postgreContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	postgreEndpoint, err := postgreContainer.Endpoint(context.Background(), "")
	if err != nil {
		t.Error(err)
	}

	defer func(postgreContainer testcontainers.Container, ctx context.Context) {
		terminationErr := postgreContainer.Terminate(ctx)
		if terminationErr != nil {
			t.Error(terminationErr)
		}
	}(postgreContainer, ctx)

	log.Println("Bootstrapping...")

	myLogger := logger.NewLogger(logger.WithPretty(true))

	// Application
	app, err := api.InitializeAppForTesting(myLogger)
	if err != nil {
		t.Fatalf("failed initializing app: %s", err)
	}
	app.Config.DatabaseUrl = fmt.Sprintf("postgresql://postgres:example@%s/db", postgreEndpoint)

	testInterrupt := make(chan error)

	closed, err := api.Bootstrap(app, myLogger, api.WithInterrupt(signaling.Forward(testInterrupt)))
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		var getImagesErr error
		defer func() {
			testInterrupt <- getImagesErr
		}()

		t.Log("Fetching images...")
		result, getImagesErr := getImages()
		if getImagesErr == nil {
			t.Log("got images", result)
		}
	}()
	t.Log("waiting for test to finish...")
	<-closed

	log.Println("Ended")
}
