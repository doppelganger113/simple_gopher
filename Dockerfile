##
## Build
##

FROM golang:1.16-buster AS build

ENV CGO_ENABLED=0
ENV GOMAXPROCS=1

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/
COPY Makefile /app/
RUN go mod download

COPY *.go /app/

RUN make

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/api /api

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/api"]
