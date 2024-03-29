SHELL = /bin/bash
PROJECTNAME=api

GOBASE=$(shell pwd)
GOBIN=bin

TZ=UTC
APP_BUILD_DATE=$(shell date +%c)
GIT_COMMIT_SHA=$(shell git rev-list -1 HEAD)

default:
	go build -v -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go || exit

install:
	go mod download

production:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build \
		-ldflags="-s -w -X 'main.Version=v1.0.0' -X 'main.GitCommit=${GIT_COMMIT_SHA}' -X 'main.BuildDate=${APP_BUILD_DATE}'" \
		-o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go

start:
	go build -v -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go|| exit
	$(GOBIN)/$(PROJECTNAME) || exit

verify:
	go test ./... -race

cover:
	go test ./... -race -coverprofile="c.out" && go tool cover -func=c.out

cover-html:
	go test ./... -race -coverprofile="c.out" && go tool cover -html=c.out

vet:
	go vet ./...

check:
	staticcheck ./...

migrate:
	go build -v -o $(GOBIN)/migrate ./cmd/migrate/main.go || exit

migrate_up:
	go run ./cmd/migrate/main.go

migrate_up_step:
	go run ./cmd/migrate/main.go -steps 1

# Migrate down by step is safer!
migrate_down:
	go run ./cmd/migrate/main.go -steps -1

# Tidy up dependencies
tidy:
	go mod tidy
