package simple_gopher

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"mime/multipart"
	"simple_gopher/auth"
	"simple_gopher/image_resize"
	"simple_gopher/storage"
)

func (service *ImagesService) getMultipleSignUrls(
	ctx context.Context, authHeader string, format image_resize.ImageFormat,
) (
	*image_resize.SignedResponse,
	*image_resize.SignedResponse,
	error,
) {
	var firstRes *image_resize.SignedResponse
	var secondRes *image_resize.SignedResponse

	g := new(errgroup.Group)

	g.Go(func() error {
		result, err := service.resizeApi.FetchSignedUrl(ctx, authHeader, format)
		if err == nil {
			firstRes = result
		}
		return err
	})

	g.Go(func() error {
		result, err := service.resizeApi.FetchSignedUrl(ctx, authHeader, format)
		if err == nil {
			secondRes = result
		}
		return err
	})

	err := g.Wait()
	if err != nil {
		return nil, nil, err
	}

	return firstRes, secondRes, nil
}

func (service *ImagesService) uploadBothFiles(
	ctx context.Context,
	originalSignedUrl string,
	croppedSignedUrl string,
	format image_resize.ImageFormat,
	original *multipart.FileHeader,
	cropped *multipart.FileHeader,
) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		return service.resizeApi.UploadFile(ctx, originalSignedUrl, format, original)
	})

	g.Go(func() error {
		return service.resizeApi.UploadFile(ctx, croppedSignedUrl, format, cropped)
	})

	err := g.Wait()

	return err
}

func (service *ImagesService) UploadAndResize(
	ctx context.Context,
	authorization auth.AuthorizationDto,
	imageName string,
	format image_resize.ImageFormat,
	originalFile *multipart.FileHeader,
	croppedFile *multipart.FileHeader,
) (*storage.Image, error) {
	seoImageName := FormatForSeo(imageName)
	if seoImageName == "" {
		return nil, InvalidArgument{
			Reason: fmt.Sprintf("Invalid image name of %s", imageName),
		}
	}

	currentUser, err := service.authenticator.GetOrSyncUser(ctx, authorization)
	if err != nil {
		return nil, err
	}

	isNameTaken, err := service.imagesRepository.DoesImageExist(ctx, seoImageName)
	if err != nil {
		return nil, InvalidArgument{
			Reason: fmt.Sprintf("Image '%s' already exists", seoImageName),
		}
	}
	if isNameTaken {
		return nil, InvalidArgument{
			Reason: fmt.Sprintf("Image name: '%s' already exists, please use another", seoImageName),
		}
	}

	originalSigned, croppedSigned, err := service.getMultipleSignUrls(ctx, authorization.Header, format)
	if err != nil {
		return nil,
			fmt.Errorf("error creating multiple sign urls: %w", err)
	}

	err = service.uploadBothFiles(
		ctx,
		originalSigned.SignedUrl,
		croppedSigned.SignedUrl,
		format,
		originalFile,
		croppedFile,
	)
	if err != nil {
		return nil, fmt.Errorf("error uploading files: %w", err)
	}

	resizeRequest := image_resize.ImageResizeRequest{
		Name:             seoImageName,
		FilePath:         croppedSigned.FileName,
		OriginalFilePath: originalSigned.FileName,
	}
	res, err := service.resizeApi.Resize(ctx, authorization.Header, resizeRequest)
	if err != nil {
		return nil, fmt.Errorf("error resizing: %w", err)
	}

	newImage := storage.Image{
		Name:     res.Name,
		Format:   storage.ImageFormat(res.Format),
		Original: res.Original,
		Domain:   res.Domain,
		Path:     res.Path,
		Sizes:    convertImageSizesToStorageSizes(res.Sizes),
		AuthorId: currentUser.Id,
	}

	createdImg, err := service.imagesRepository.Create(ctx, newImage)
	if err != nil {
		return nil, fmt.Errorf("err saving new image to database: %w", err)
	}

	return createdImg, nil
}
