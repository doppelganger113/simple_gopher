SHELL = /bin/bash
PROJECTNAME=simple_gopher

GOBASE=$(shell pwd)
GOBIN=bin

default:
	go build -i -v -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go || exit

install:
	go mod download

production:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go

start:
	go build -i -v -o $(GOBIN)/$(PROJECTNAME) ./cmd/$(PROJECTNAME)/main.go|| exit
	$(GOBIN)/$(PROJECTNAME) || exit

verify:
	go test ./... -race

cover:
	go test ./... -race -coverprofile="c.out" && go tool cover -func=c.out

cover-html:
	go test ./... -race -coverprofile="c.out" && go tool cover -html=c.out

vet:
	go vet ./...

migrations:
	go build -i -v -o $(GOBIN)/migrations ./cmd/migrations/main.go || exit

migrate_up:
	go run ./cmd/migrations/main.go

migrate_up_step:
	go run ./cmd/migrations/main.go -steps 1

# Migrate down by step is safer!
migrate_down:
	go run ./cmd/migrations/main.go -steps -1

# Tidy up dependencies
tidy:
	go mod tidy
