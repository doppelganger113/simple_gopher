package http_transport

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"simple_gopher/storage"
	"testing"
)

func TestFetchImage(t *testing.T) {
	config := Config{}
	handlers := Handlers{
		ImagesHandler: ImagesHandlerMock{},
		Authenticator: AuthenticatorMock{},
	}
	ctx := context.Background()
	server, err := NewServer(config, handlers)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = server.Shutdown(ctx); err != nil {
			t.Error(err)
		}
	}()

	testServer := httptest.NewServer(server.router)
	defer testServer.Close()

	res, err := http.Get(testServer.URL + "/api/v1/images?size=10&page=2&order=ASC")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %d", res.StatusCode)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			t.Error(err)
		}
	}()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var images storage.ImageList
	if err = json.Unmarshal(data, &images); err != nil {
		t.Fatal(err)
	}

	expectedImages := storage.ImageList{
		{
			Id:       "3c47d736-6c4e-4a1c-a04b-3744cc30b263",
			Name:     "my-image-1",
			Format:   "jpg",
			Original: "images/my-unique-image.png",
			Domain:   "https://random.cloudfront.net",
			Path:     "images",
			Sizes: storage.ImageSizes{
				Original: storage.Dimensions{
					Width:  688,
					Height: 516,
				},
				Xs: &storage.Dimensions{
					Width:  100,
					Height: 75,
				},
				S: &storage.Dimensions{
					Width:  300,
					Height: 225,
				},
				M: &storage.Dimensions{
					Width:  500,
					Height: 375,
				},
			},
			CreatedAt: nil,
			UpdatedAt: nil,
		},
	}

	if !reflect.DeepEqual(images, expectedImages) {
		t.Fatal("result does not match expected images")
	}
}

//func (s *MySuite) TestFetchImages() {
//	repoMock := storage.ImageRepoMock{}
//
//	server, err := NewServer(app)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	defer func() {
//		_ = app.Shutdown(context.Background())
//	}()
//
//	testServer := httptest.NewServer(server.router)
//	defer testServer.Close()
//
//	res, err := http.Get(testServer.URL + "/api/v1/images")
//	assert.Nil(s.T(), err)
//	assert.Equal(s.T(), res.StatusCode, 200)
//
//	defer func() {
//		_ = res.Body.Close()
//	}()
//	data, err := ioutil.ReadAll(res.Body)
//	assert.Nil(s.T(), err)
//
//	repoMock.AssertCalled(s.T(), "Get", mock.Anything, 20, 0, storage.OrderDescending)
//
//	var images storage.ImageList
//	err = json.Unmarshal(data, &images)
//	assert.Nil(s.T(), err)
//
//	expectedImages := storage.ImageList{
//		{
//			Id:       "3c47d736-6c4e-4a1c-a04b-3744cc30b263",
//			Name:     "my-image-1",
//			Format:   "jpg",
//			Original: "images/my-unique-image.png",
//			Domain:   "https://random.cloudfront.net",
//			Path:     "images",
//			Sizes: storage.ImageSizes{
//				Original: storage.Dimensions{
//					Width:  688,
//					Height: 516,
//				},
//				Xs: &storage.Dimensions{
//					Width:  100,
//					Height: 75,
//				},
//				S: &storage.Dimensions{
//					Width:  300,
//					Height: 225,
//				},
//				M: &storage.Dimensions{
//					Width:  500,
//					Height: 375,
//				},
//			},
//			CreatedAt: nil,
//			UpdatedAt: nil,
//		},
//	}
//
//	assert.Equal(s.T(), images, expectedImages)
//}
//
//func (s *MySuite) TestUploadFile() {
//	repoMock := new(storage.ImageRepoMock)
//	repoMock.
//		On("Create", mock.Anything, storage.Image{
//			Name:     "my-image-1",
//			Format:   "png",
//			Original: "images/my-unique-image.png",
//			Domain:   "https://random.cloudfront.net",
//			Path:     "images",
//			Sizes: storage.ImageSizes{
//				Original: storage.Dimensions{
//					Width:  200,
//					Height: 150,
//				},
//				Xs: &storage.Dimensions{
//					Width:  100,
//					Height: 75,
//				},
//				S: &storage.Dimensions{
//					Width:  150,
//					Height: 100,
//				},
//			},
//			AuthorId: "user-uuid",
//		}).
//		Return(&storage.Image{
//			Id:       "3c47d736-6c4e-4a1c-a04b-3744cc30b263",
//			Name:     "my-image-1",
//			Format:   "png",
//			Original: "images/my-unique-image.png",
//			Domain:   "https://random.cloudfront.net",
//			Path:     "images",
//			AuthorId: "user-uuid",
//			Sizes: storage.ImageSizes{
//				Original: storage.Dimensions{
//					Width:  200,
//					Height: 150,
//				},
//				Xs: &storage.Dimensions{
//					Width:  100,
//					Height: 75,
//				},
//				S: &storage.Dimensions{
//					Width:  150,
//					Height: 100,
//				},
//			},
//		}, nil)
//	repoMock.
//		On("DoesImageExist", mock.Anything, "my-dummy").
//		Return(false, nil)
//
//	authenticationMock := new(auth.AuthMock)
//	authenticationMock.
//		On("IsTokenValid", mock.Anything, bearerTokenMock, storage.AuthRoleAdmin).
//		Return(true, "username-uuid", nil)
//	authenticationMock.On("Shutdown").Return(nil)
//	authenticationMock.
//		On("GetOrSyncUser", mock.Anything, auth.AuthorizationDto{
//			Header:   authHeaderMock,
//			Username: "username-uuid",
//			Role:     auth.RoleAdmin,
//		}).
//		Return(&storage.User{
//			Id:          "user-uuid",
//			Email:       "john@gmail.com",
//			CreatedAt:   time.Time{},
//			Role:        storage.AuthRoleAdmin,
//			CogUsername: "cog-username",
//			CogSub:      "cog-sub",
//			CogName:     "cog-name",
//			Disabled:    false,
//		}, nil)
//
//	resizeApiMock := new(image_resize.Mock)
//	resizeApiMock.
//		On("FetchSignedUrl", mock.Anything, authHeaderMock, image_resize.PngFormat).
//		Return(
//			&image_resize.SignedResponse{
//				SignedUrl: "test-signed-url",
//				FileName:  "random-number",
//			},
//			nil,
//		)
//	resizeApiMock.
//		On("UploadFile", mock.Anything, "test-signed-url", image_resize.PngFormat, mock.Anything).
//		Return(nil)
//	resizeApiMock.
//		On("Resize", mock.Anything, authHeaderMock, image_resize.ImageResizeRequest{
//			Name:             "my-dummy",
//			FilePath:         "random-number",
//			OriginalFilePath: "random-number",
//		}).
//		Return(
//			&image_resize.ImageResizeResponse{
//				Format:   "png",
//				Name:     "my-image-1",
//				Original: "images/my-unique-image.png",
//				Domain:   "https://random.cloudfront.net",
//				Path:     "images",
//				Sizes: image_resize.ImageSizes{
//					Original: image_resize.Dimensions{
//						Width:  200,
//						Height: 150,
//					},
//					Xs: &image_resize.Dimensions{
//						Width:  100,
//						Height: 75,
//					},
//					S: &image_resize.Dimensions{
//						Width:  150,
//						Height: 100,
//					},
//				},
//			},
//			nil,
//		)
//
//	app := simple_gopher.NewApp(
//		simple_gopher.Config{},
//		new(storage.Mock),
//		authenticationMock,
//		simple_gopher.NewImagesService(
//			resizeApiMock,
//			repoMock,
//			authenticationMock,
//		),
//	)
//	server, err := NewServer(app)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	defer func() {
//		_ = app.Shutdown(context.Background())
//	}()
//
//	testServer := httptest.NewServer(server.router)
//	defer testServer.Close()
//
//	croppedFile := test.FormFilesParams{
//		FilePath:  "../test/dummy_image.png",
//		FileField: "croppedFile",
//		FileName:  "my-dummy-cropped.png",
//	}
//	originalFile := test.FormFilesParams{
//		FilePath:  "../test/dummy_image.png",
//		FileField: "originalFile",
//		FileName:  "my-dummy.png",
//	}
//
//	formData, contentType, err := test.CreateMultipartFormData(
//		[]test.FormFilesParams{croppedFile, originalFile},
//		test.FormParams{
//			"name":   "my-dummy",
//			"format": string(image_resize.PngFormat),
//		},
//	)
//	assert.Nil(s.T(), err)
//
//	req, err := http.NewRequest("POST", testServer.URL+"/api/v1/images", formData)
//	assert.Nil(s.T(), err)
//	req.Header.Set("Authorization", authHeaderMock)
//	req.Header.Set("Content-Type", contentType)
//
//	client := new(http.Client)
//	resp, err := client.Do(req)
//	assert.Nil(s.T(), err)
//	defer func() {
//		_ = resp.Body.Close()
//	}()
//
//	assert.Equal(s.T(), resp.StatusCode, 201)
//
//	response, err := ioutil.ReadAll(resp.Body)
//	assert.Nil(s.T(), err)
//
//	var createdImage storage.Image
//	err = json.Unmarshal(response, &createdImage)
//	assert.Nil(s.T(), err)
//
//	expected := storage.Image{
//		Id:       "3c47d736-6c4e-4a1c-a04b-3744cc30b263",
//		Name:     "my-image-1",
//		Format:   "png",
//		Original: "images/my-unique-image.png",
//		Domain:   "https://random.cloudfront.net",
//		Path:     "images",
//		AuthorId: "user-uuid",
//		Sizes: storage.ImageSizes{
//			Original: storage.Dimensions{
//				Width:  200,
//				Height: 150,
//			},
//			Xs: &storage.Dimensions{
//				Width:  100,
//				Height: 75,
//			},
//			S: &storage.Dimensions{
//				Width:  150,
//				Height: 100,
//			},
//		},
//	}
//
//	assert.Equal(s.T(), expected, createdImage)
//}
