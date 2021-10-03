package simple_gopher

import (
	"context"
	"github.com/google/uuid"
	"log"
	"simple_gopher/auth"
	"simple_gopher/image"
	"simple_gopher/storage"
)

func (service *ImagesService) DeleteOne(
	ctx context.Context,
	auth auth.AuthorizationDto,
	imageId string,
) error {
	parsedId, err := uuid.Parse(imageId)
	if err != nil {
		return InvalidArgument{Reason: "invalid uui"}
	}

	user, err := service.authenticator.GetOrSyncUser(ctx, auth)
	if err != nil {
		return err
	}
	if user.Role != storage.AuthRoleAdmin {
		return Forbidden{}
	}

	img, err := service.imagesRepository.GetOne(ctx, imageId)
	if err != nil {
		return err
	}

	deleteRequest := image.DeleteRequest{
		Name:       img.Name,
		Format:     image.Format(img.Format),
		Dimensions: convertStorageSizesToDimensions(img.Sizes),
	}
	if err = service.resizeApi.Delete(ctx, auth.Header, deleteRequest); err != nil {
		log.Println("failed deleting image " + imageId)
		return err
	}
	if err = service.resizeApi.Invalidate(ctx, auth.Header, deleteRequest); err != nil {
		log.Println("failed invalidating image " + imageId + " " + err.Error())
	}

	return service.imagesRepository.DeleteOne(ctx, parsedId.String())
}
