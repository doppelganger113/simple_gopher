# Architecture

## Introduction

Often frameworks lock us in instead of providing us with a bootstrap to get started quickly and develop rapidly, we end
up needing to learn the framework, how it works, follow its structure guideline and if we hit a wall then we need to dig
deep into the framework, the framework needs to be updated and maintained and becomes a huge bag of things that we
constantly have to carry with us and tend to, this becomes burden over time. Golang already provides us with a powerful
standard library that is constantly updated and which we can utilize to pretty much any extent our needs might require.

## Hexagonal architecture

Hexagonal architecture is an approach in which we extract code/logic which communicates with the "outside world",
examples of this are database calls, network calls like HTTP RPC SFTP AMQP, UI, etc, everything that is not our
application itself. Main reasons for this is decoupling that eases maintenance and testability through mocks.

### Project directory structure

#### Storage/Database and HTTP API clients

Directory structure is a bit of domain driven design and package level design that Go has, which can be seen in this
example. Let's start with a single example of **storage/database**.

```text
simple_gopher/
├─ README.md
├─ Makefile
├─ storage/
│  ├─ storage.go
│  ├─ storage_mock.go
│  ├─ image.go
│  ├─ image_repository.go
│  ├─ image_repository_mock.go
│  ├─ postgresql/
│  │  ├─ database.go
│  │  ├─ image_repository.go
│  │  ├─ image_repository_test.go
```

Storage acts as a self-contained package which has in it's first layer the interface for the database and repositories
and models that are 1:1 mapped against the database table/document. Second layer is the actual implementation of it. It
is done this way to avoid circular dependency, to allow us to mock the actual database calls or easily replace the
underlying implementation.

Implementation can have tests (`storage/postgresql/image_repository_test.go`) that will test the actual DB calls if we
want to perform **integration testing**. We can use a flag in these integration tests to separate execution of unit and
integration tests.

The `storage/storage_mock.go` and `storage/image_repository_mock.go` are used to prepare mock implementation that we can
use in tests and override when required so. Example of this would look like:

```go
package storage

import (
	"context"
	"errors"
)

var DuplicateErr = errors.New("duplicate, already exists")

type NotFound struct {
	Msg string
}

func (nf NotFound) Error() string {
	return nf.Msg
}

type Image struct {
	Id   string
	Name string
}

type ImageRepository interface {
	GetOneById(ctx context.Context, id string) (Image, error)
	Create(ctx context.Context, image Image) (Image, error)
}

type ImageRepositoryMock struct {
}

func (repo ImageRepositoryMock) GetOneById(_ context.Context, id string) (Image, error) {
	if id != "1" {
		return Image{}, NotFound{Msg: "Could not find image with id: " + id}
	}
	return Image{Id: 0, Name: "my-image"}, nil
}

func (repo ImageRepositoryMock) Create(_ context.Context, image Image) (Image, error) {
	if image.Id == "1" {
		return Image{}, DuplicateErr
	}
	return image, nil
}
```
Now we can re-use this mock and override it in tests however we need
```go
package simple_gopher

import (
	"context"
	"testing"
)

type HelloService struct {
	repo ImageRepository
}

func NewHelloService(repo ImageRepository) HelloService {
	return HelloService{repo: repo}
}

func (hello HelloService) CreateImage(ctx context.Context, image storage.Image) (storage.Image, error) {
	return hello.repo.Create(ctx, image)
}

type ImageRepositoryMock struct {
	storage.ImageRepositoryMock // we use composition to "inherit" methods
}

// We override the existing Create method with a new one, though note that you can update
// the existing mock to have complex logic which can be re-used

func (repo ImageRepositoryMock) Create(_ context.Context, image Image) (Image, error) {
	return image, nil
}

func Test_SomeFunc(t *testing.T) {
	repoMock := ImageRepositoryMock{}
	expected := storage.Image{Id: 3, Name: "new-img"}

	img, err := NewHelloService(repoMock).
		CreateImage(context.Background(), expected)
	if err != nil {
		t.Fataln("failed creating image: %v", err)
	}

	if img != expected {
		t.Fatalf("images don't match, expected: %v got: %v", expected, img)
	}
}
```

This strategy applies to all "outer" logic like with `image_resize` package which wraps HTTP cli that communicates with
image resize service.

**Note:** Prepare to copy some structures and types in this approach, in order to avoid circular dependency, side
effects and decoupling some stuff will need a duplicate in another package. Intention is to keep the packages as
decoupled as possible and think of it in a way that you can fully copy the package into another project without pulling
any dependencies with it.

#### Transport layer structure

In this case we are using http as the transport layer in `http_transport/` package and directory, but it is made in such
a way that you can easily plug it out and replace with gRPC for example. Transport layer libraries are a simple glue for
our application, they connect input and handle output if any.

In order to achieve this we have to write a bit more code and by that I mean we need to create interfaces.

```go
package http_transport

type ImagesHandler interface {
	Get(ctx context.Context, limit, offset int, order storage.Order) (storage.ImageList, error)
}
type Handlers struct {
	ImagesHandler ImagesHandler
}

func NewServer(config Config, handlers Handlers) (*Server, error) {
	
	// ...some initial setup middleware, configuration, etc.
	// register handlers for routes
	r.Route("/api/v1/images", ImagesRouter(handlers.ImagesHandler))
	
	return &Server{}, nil
}
```
This way we don't care how our app works, we just map the REST API endpoints to those methods. When testing we can
easily test with interfaces if query params are passed correctly and if validation is done correctly, this way we 
achieve decoupling and make testing easier.

#### CMD

This directory serves as main entrypoint for constructing and building application/s. It follows the following structure
where each app gets its own directory with single `main.go` file.

```text
simple_gopher/
├─ cmd/
│  ├─ simple_gopher/
│  │  ├─ main.go
│  ├─ migrate/
│  │  ├─ main.go
```
These `main.go` files are simple and are used only to bootstrap the application


#### Core application

Core application is usually placed in a directory with a name that corresponds to the service name or what it does. 
This is on order to follow the package semantics, so when we import it we know what it is based on it's name
`simple_gopher.NewApp()` instead of `src.NewApp()` or `core.NewApp()`.

```text
simple_gopher/
├─ simple_gopher/
│  │  ├─ app.go
```

In this `app.go` file we construct our application by manual dependency injection initializing implementations of 
storage, image resize and other "ports/adapters". We can add methods that are initialization or shutdown of our 
application as well.

```go
package simple_gopher

func CreateApp(config Config) (*App, error) {
	authConfig := NewAuthConfig(config)

	// Storage
	db := postgresql.NewDatabase()
	imageRepository := postgresql.NewImageRepository(db)
	userRepository := postgresql.NewUserRepo(db)

	// Services
	imageResizeConfig := image_resize_api.Config{ImagesApiDomain: config.ImagesApiDomain}
	resizeApi := image_resize_api.NewResizeApi(imageResizeConfig)
	authenticator := cognito.NewCognitoAuthService(authConfig, userRepository)
	imagesService := NewImagesService(resizeApi, imageRepository, authenticator)

	// App
	app := NewApp(
		config, db, authenticator, imagesService,
	)

	return app, nil
}
```

Basically, this directory should contain services which have our business logic, you can choose 2 ways of organizing 
files:

1. Have a single file for service
    ```text
    images.go // Contains all methods like get, create, update, delete
    ```
2. Split methods in multiple files to reduce the file size, use following naming approach in order to group them
   ```text
   images.go
   images_get.go
   images_create.go
   images_update.go
   images_delete.go
   ```
Don't be afraid to have many files in one directory.

Tests in this directory would be the most important tests as they test your business logic, while the others tested if
integration works. You will use mocks of "ports/adapters"

##### Core application integration testing

If needed, you can create some test utility functions that on a given flag construct the app with real implementation
of packages instead of using mocks to test the services fully integrated or partially.

### TODO

Update migrations and add some documentation regarding it.