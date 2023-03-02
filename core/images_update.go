package core

import (
	"api/auth"
	"api/core/exception"
	"api/image"
	"api/storage"
	"context"
	"fmt"
	"mime/multipart"
)

func (service *ImagesService) Update(
	ctx context.Context,
	imageId string,
	authorization auth.AuthorizationDto,
	imageName string,
	format image.Format,
	originalFile *multipart.FileHeader,
	croppedFile *multipart.FileHeader,
) (storage.Image, error) {
	isFileUpload := string(format) != "" && originalFile != nil && croppedFile != nil
	if !isFileUpload && imageName == "" {
		return storage.Image{}, exception.InvalidArgument{
			Reason: "Expected at least a name change or a file change, got all empty",
		}
	}

	_, err := service.authenticator.GetOrSyncUser(ctx, authorization)
	if err != nil {
		return storage.Image{}, err
	}

	img, err := service.imagesRepository.GetOne(ctx, imageId)
	if err != nil {
		return storage.Image{}, err
	}

	if isFileUpload && imageName != "" {
		return service.updateImageAndName(
			ctx, authorization.Header, imageName, format, img, originalFile, croppedFile,
		)
	} else if imageName == "" {
		return service.updateImageOnly(
			ctx, authorization, imageId, format, originalFile, croppedFile,
		)
	}

	_, err = service.updateNameOnly(ctx, authorization.Header, img, imageName)
	if err != nil {
		return storage.Image{}, err
	}

	return service.imagesRepository.SetNameById(ctx, imageId, imageName)
}

func (service *ImagesService) updateNameOnly(
	ctx context.Context,
	authHeader string,
	img storage.Image,
	newImageName string,
) (image.ResizeResponse, error) {
	request := image.RenameRequest{
		Name:    img.Name,
		NewName: newImageName,
		Format:  image.Format(img.Format),
		SizeMap: fromStorageImageSizesToImageSizes(img.Sizes),
	}

	response, err := service.resizeApi.Rename(ctx, authHeader, request)
	if err != nil {
		return image.ResizeResponse{}, err
	}

	return response, nil
}

func (service *ImagesService) updateImageOnly(
	ctx context.Context,
	authDto auth.AuthorizationDto,
	imageId string,
	format image.Format,
	originalFile *multipart.FileHeader,
	croppedFile *multipart.FileHeader,
) (storage.Image, error) {
	img, err := service.imagesRepository.GetOne(ctx, imageId)
	if err != nil {
		return storage.Image{}, err
	}

	originalSignedUrl, croppedSignedUrl, err := service.
		getMultipleSignUrls(ctx, authDto.Header, format)
	if err != nil {
		return storage.Image{}, err
	}

	request := image.DeleteRequest{
		Name:       img.Name,
		Format:     format,
		Dimensions: convertStorageSizesToDimensions(img.Sizes),
	}
	err = service.resizeApi.Delete(ctx, authDto.Header, request)
	if err != nil {
		return storage.Image{}, err
	}

	err = service.uploadBothFiles(
		ctx,
		originalSignedUrl.SignedUrl,
		croppedSignedUrl.SignedUrl,
		format,
		originalFile,
		croppedFile,
	)
	if err != nil {
		return storage.Image{}, err
	}

	imageResizeRequest := image.ResizeRequest{
		Name:             img.Name,
		FilePath:         croppedSignedUrl.FileName,
		OriginalFilePath: originalSignedUrl.FileName,
	}
	_, err = service.resizeApi.Resize(ctx, authDto.Header, imageResizeRequest)
	if err != nil {
		return storage.Image{}, err
	}

	return img, nil
}

func (service *ImagesService) updateImageAndName(
	ctx context.Context,
	authHeader string,
	imageName string,
	format image.Format,
	img storage.Image,
	originalFile *multipart.FileHeader,
	croppedFile *multipart.FileHeader,
) (storage.Image, error) {
	seoImageName := FormatForSeo(imageName)
	if seoImageName == "" {
		return storage.Image{}, exception.InvalidArgument{
			Reason: fmt.Sprintf("Invalid image name of %s", imageName),
		}
	}

	original, cropped, err := service.getMultipleSignUrls(ctx, authHeader, format)
	if err != nil {
		return storage.Image{}, err
	}

	if err = service.uploadBothFiles(
		ctx, original.SignedUrl, cropped.SignedUrl, format, originalFile, croppedFile,
	); err != nil {
		return storage.Image{}, err
	}

	request := image.DeleteRequest{
		Name:       imageName,
		Format:     format,
		Dimensions: convertStorageSizesToDimensions(img.Sizes),
	}
	if err = service.resizeApi.Delete(ctx, authHeader, request); err != nil {
		return storage.Image{}, err
	}

	resizeRequest := image.ResizeRequest{
		Name:             imageName,
		FilePath:         cropped.FileName,
		OriginalFilePath: original.FileName,
	}
	res, err := service.resizeApi.Resize(ctx, authHeader, resizeRequest)
	if err != nil {
		return storage.Image{}, err
	}

	newImage := storage.Image{
		Name:     res.Name,
		Format:   storage.ImageFormat(res.Format),
		Original: res.Original,
		Domain:   res.Domain,
		Path:     res.Path,
		Sizes:    convertImageSizesToStorageSizes(res.Sizes),
		AuthorId: img.AuthorId,
	}
	if err = service.imagesRepository.UpdateOne(ctx, newImage); err != nil {
		return storage.Image{}, err
	}

	return newImage, nil
}
