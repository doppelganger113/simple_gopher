##
## Build
##

FROM golang:1.20-buster AS build

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
COPY Makefile /app/

RUN go mod download
RUN go mod verify

COPY *.go /app/

RUN make production

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /core/bin/api /api

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/api"]
