# Simple Gopher

**Introduction**

You'd probably want to check out [ARCHITECTURE.md](./ARCHITECTURE.md) file to learn more on the architecture of the 
project layout.

Partial example of a service using hexagonal architecture approach with Go package oriented design. This serves as an
example how to design Go services without a framework, high decoupling and easy testability. The example isn't 100%
completed or fully covered, but good enough example on how it should look like. It will probably gain more polishing
as there is always stuff to be added or updated, but for now, lets start with this.

Bellow is an example of how the README.md file for the service should look like.

<div align="center">
    <img src="./assets/go_logo.png" align="center" width="200" alt="Go" />
</div>

## Description

Won't get too much into the service itself, it's for you to discover and isn't too important, but to keep it shot, 
the Simple Gopher is a web API service that integrates with image resize service for resizing images and storing that 
data and integrates with AWS Cognito for authentication.

Table of contents
=================

<!--ts-->

* [Configuration](#configuration)
* [Developing](#developing)
    * [Development requirements](#development-requirements)
    * [Running locally](#running-locally)
    * [Dependency management](#running-locally)
* [Open Api 3 documentation](#open-api-3-documentation)
* [Testing](#testing)
    * [Unit testing](#unit-tests)
    * [Integration testing](#integration-tests)
    * [Test configuration](#test-configuration)
* [Migrations](#migrations)
    * [Migration requirements](#migration-requirements)
    * [Running migrations](#running-migrations)
    * [Migrations in CI/CD](#migrations-in-cicd)
* [CI/CD](#cicd)
    * [Deploying to CI/CD](#deploying-to-cicd)
* [Troubleshooting](#troubleshooting)
* [Helpful materials](#helpful-materials)

<!--te-->

## Configuration

Configuration of the application is done through the environment variables, which are following:

**Required** environment variables:

| Environment variable name         | Required  | Description                                 |
| -------------                     | --------- | ------------------------------------------- |
| PORT                              | Optional  | Default value is 3000                       |
| DEBUG_ROUTES                      | Optional  | Default value is false, set to `true` for local endpoint debugging                       |
| AWS_REGION                        | Required  | Example: `eu-central-1`                                            |
| AWS_USER_POOL_ID                  | Required  | Example: `eu-central-1_somenumber`                                            |
| AWS_ACCESS_KEY_ID                 | Required  | AWS IAM key id used for API to talk with AWS services                                            |
| AWS_SECRET_ACCESS_KEY             | Required  | Secret for the above key                                            |
| DATABASE_URL                      | Required  | Example `postgresql://postgres:example@localhost/db?sslmode=disable`                                            |
| IMAGES_API_DOMAIN                 | Required  | Endpoint for the image service API                                            |
| CORS_ALLOW_ORIGINS                | Required  | List of origins to allow CORS in format: `first.com, second.com, etc.com`                                            |
| SQS_POST_AUTH_URL                 | Required  | Url of the SQS queue                                            |
| SQS_POST_AUTH_INTERVAL_SEC        | Optional  | Interval in which the API will pool the queue for user registration events. Default value is `600`                                            |
| SQS_POST_AUTH_CONSUMER_DISABLED   | Optional  | Default value false, set value to `true` to turn off in modes like local development to avoid messing with production                                            |
| BASIC_AUTH_REALM                  | Optional  | Name of the realm for authentication, default is Forbidden |
| BASIC_AUTH_USERNAME               | Optional  | Username used for basic authentication |
| BASIC_AUTH_PASSWORD               | Optional  | Password used for basic authentication |
| OAUTH2_AUTHORIZATION_CODE_URL     | Optional  | Url for OAuth2 authentication in format `https://your-domain.auth.eu-central-1.amazoncognito.com/login?response_type=code&client_id=<your-client-id>&redirect_uri=<your-redirect-uri>` |
| OAUTH2_TOKEN_URL                  | Optional  | Url for OAuth2 token retrieval in format `https://your-domain.auth.eu-central-1.amazoncognito.com/oauth2/token` |
| DOMAIN                            | Optional  | Name of the domain the app is being served from, like `localhost:3000` or `https://your-domain.herokuapp.com` |

## Developing

### Development requirements

- Go v1.17
- PostgreSQL
- Docker and docker-compose

You will need to create a new aws cli profile locally and use it with credentials for this API. To create it
execute `aws configure --profile your-profile` then `export AWS_PROFILE=your-profile`

### Running locally

1. Set the environment variables. It's also recommended storing them in `.env` file which is ignored for ease of
   management. Export variables in each row like `export MY_VARIABLE=1234`, then load the env variables
   with `source .env`
2. Ensure that you start the database and all the other required services by running:
   `docker-compose up -d`
3. After setting the environment variables, execute the `make start` command to build and start the server

### Dependency management

- **Remove dependency** by removing all occurrences of the library in imports and execute:
  ```bash
  go mod tidy
  ```

## Open Api 3 documentation

Golang implementation of [OpenApi3 specification](https://swagger.io/docs/specification/basic-structure/) aka Swagger
through dynamic configuration with Swagger UI.

We use Swagger-UI with some small changes in order for it to fetch changes from our API where we can set the redirection
url for OAuth2.

**Dependencies**

- github.com/getkin/kin-openapi/openapi3
- github.com/getkin/kin-openapi/openapi3gen

Documentation is added/updated in `src/api/http-transport/openapi3.go`, served via
`src/api/http-transport/handler_openapi3.go` with swagger UI located in `src/api/http-transport/docs`.

## Testing

### Test configuration

Ensure that you set the `DATABASE_TEST_URL` environment variable, reason for new one is for safety reasons.

```bash
export DATABASE_TEST_URL=postgresql://postgres:example@localhost/db?sslmode=disable

# Also, disable the pooling, we don't need it for testing
export SQS_POST_AUTH_CONSUMER_DISABLED=true
```

### Unit testing

- `make verify`

### Integration testing

Unit testing **with** integration testing

- `export TEST_INTEGRATION=true`
- `make verify`

### Vetting

Go provides a great tool for checking out code and detecting possible bugs, think of it as a linter. To run it execute:

```bash
make vet
```

### Testing in CI/CD

Ensure that you run the migrations before running end-to-end tests.
`TBD: create separate e2e test execution`

## Migrations

### Migration requirements

To run the migrations, you must set the database connection url through environment variable:

```bash
export DATABASE_URL=postgresql://postgres:example@localhost/db?sslmode=disable
```

Note that the `sslmode=disable` is for local development only, for production you should use ssl encryption.

### Running migrations

**BE CAREFUL WHEN REVERTING!**
Migrations are done with Go library [golang-migrate/migrate](https://github.com/golang-migrate/migrate) and are executed
with following commands:

- `make migrate_up` to update the schema to the latest
- `make migrate_up_step` to update the migration one step up
- `make migrate_down` to revert the schema one step down

It's paramount to follow [best practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md) to ensure
you don't break anything.

### Migrations in CI/CD

For CI/CD ensure that you always run migrations before deploying your new app, like a pre-run action. For rollbacks,
also ensure that you first close the application before running the migration down command. As always, for specific
cases you will need to pay attention if you need code that supports both versions, so it doesn't crash if the change is
drastic.

## CI/CD

CI/CD is currently on the Heroku and additional options that were added for it are located in `go.mod` file as:

1. Build phase

```text
// +heroku goVersion 1.17
// +heroku install ./cmd/...
```

2. Run phase in `Procfile`

They specify the Go version and build location. Builds will end up in `bin/` directory. Other thing to note is that
the `bin/go-post-compile` and `bin/go-pre-compile` will execute if they exist, so use them for pre and post actions.

### Deploying to CI/CD

It is recommended to execute both unit and integration tests, as well as building the app before pushing to git.

```bash
# Test
export TEST_INTEGRATION=true
make test
make
```

More about Heroku options can be found at the
[Go Buildpack](https://github.com/heroku/heroku-buildpack-go#prepost-compile-hooks).

## Troubleshooting

- Error when building: `open /usr/local/go/pkg/darwin_amd64/runtime/cgo.a: permission denied`
   ```bash
   sudo chmod -R 777 /usr/local/go
   ```
- How to clear test cache? Execute `go clean -testcache`

## Helpful materials

- [Go cheat sheet](https://devhints.io/go)
- [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests/)
- [Go by example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Makefile tutorial](https://makefiletutorial.com/)

## Advanced materials

- [Network Programming with Go](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/index.html)
