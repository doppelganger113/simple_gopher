package http_server

import (
	"api/auth"
	"api/core"
	"api/core/exception"
	"api/http_server/authenticator"
	"api/http_server/http_util"
	"api/http_server/middleware"
	"api/http_server/middleware/keys"
	"api/image"
	"api/storage"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
)

const maxBodyLimitBytes = 30 * 1024 * 1024 // 20MB

type ImageHandler struct {
	http_util.RequestHandler
	imagesService *core.ImagesService
	logger        *zerolog.Logger
	authenticator authenticator.Authenticator
}

func NewImageHandler(
	logger *zerolog.Logger,
	authenticator authenticator.Authenticator,
	service *core.ImagesService,
) *ImageHandler {
	handler := http_util.NewRequestHandler(logger)

	return &ImageHandler{
		handler,
		service,
		logger,
		authenticator,
	}
}

func (h ImageHandler) CreateRouter() func(router chi.Router) {
	isAdmin := middleware.Authorize(h.logger, h.authenticator, auth.RoleAdmin)

	return func(r chi.Router) {
		r.Get("/{imageId}", h.Handle(h.fetchImage))
		r.Get("/", h.Handle(h.fetchImages))
		r.With(isAdmin).Post("/upload", h.Handle(h.addImage))
		r.With(isAdmin).Patch("/{imageId}", h.Handle(h.updateImage))
		r.With(isAdmin).Delete("/{imageId}", h.Handle(h.deleteOne))
	}
}

func (h ImageHandler) fetchImage(ctx context.Context, req *http.Request) (*http_util.Response, error) {
	imageId := chi.URLParam(req, "imageId")
	img, err := h.imagesService.GetOne(ctx, imageId)
	if err != nil {
		return nil, err
	}

	return http_util.NewResponse(img), nil
}

func (h ImageHandler) fetchImages(ctx context.Context, req *http.Request) (*http_util.Response, error) {
	page := http_util.ToUint(req.URL.Query().Get("page"))
	size := http_util.ToUint(req.URL.Query().Get("size"))

	order := storage.ToOrderOr(req.URL.Query().Get("order"), storage.OrderDescending)
	limit, offset := storage.PagingToLimitOffset(page, size)

	imageList, err := h.imagesService.Get(ctx, limit, offset, order)
	if err != nil {
		return nil, err
	}

	return http_util.NewResponse(imageList), nil
}

type UploadImageDto struct {
	Name   string
	Format image.Format
}

func (dto UploadImageDto) validate() error {
	if len(dto.Name) < 5 || len(dto.Name) > 200 {
		return exception.InvalidArgument{
			Reason: "Name should be between 5 and 250 characters",
		}
	}

	if !dto.Format.IsSupported() {
		return exception.InvalidArgument{
			Reason: fmt.Sprintf("Unsupported format %s", dto.Format),
		}
	}

	return nil
}

func (h ImageHandler) addImage(ctx context.Context, req *http.Request) (*http_util.Response, error) {
	err := req.ParseMultipartForm(maxBodyLimitBytes)
	if err != nil {
		return nil, http_util.NewFailureResponse("failed parsing multipart form data")
	}

	_, originalFileHeader, err := req.FormFile("originalFile")
	if err != nil {
		return nil, http_util.NewFailureResponse("missing originalFile")
	}
	_, croppedFileHeader, err := req.FormFile("croppedFile")
	if err != nil {
		return nil, http_util.NewFailureResponse("missing croppedFile")
	}

	data := &UploadImageDto{}
	data.Name = req.PostFormValue("name")
	data.Format = image.Format(req.PostFormValue("format"))
	if err = data.validate(); err != nil {
		return nil, http_util.NewFailureResponse(err.Error())
	}

	authorization, err := auth.ExtractAuthorizationDto(ctx, keys.UserAuthDtoKey)
	if err != nil {
		return nil, err
	}

	img, err := h.imagesService.UploadAndResize(
		ctx,
		authorization,
		data.Name,
		data.Format,
		originalFileHeader,
		croppedFileHeader,
	)
	if err != nil {
		return nil, err
	}

	return http_util.NewResponse(img).WithStatus(http.StatusCreated), nil
}

func (h ImageHandler) updateImage(ctx context.Context, req *http.Request) (*http_util.Response, error) {
	err := req.ParseMultipartForm(maxBodyLimitBytes)
	if err != nil {
		return nil, http_util.NewFailureResponse("failed parsing multipart form data")
	}

	_, originalFileHeader, err := req.FormFile("originalFile")
	if err != nil {
		return nil, http_util.NewFailureResponse("missing originalFile")
	}
	_, croppedFileHeader, err := req.FormFile("croppedFile")
	if err != nil {
		return nil, http_util.NewFailureResponse("missing croppedFile")
	}

	data := &UploadImageDto{}
	data.Name = req.PostFormValue("name")
	data.Format = image.Format(req.PostFormValue("format"))
	if err = data.validate(); err != nil {
		return nil, http_util.NewFailureResponse(err.Error())
	}

	authorization, err := auth.ExtractAuthorizationDto(ctx, keys.UserAuthDtoKey)
	if err != nil {
		return nil, err
	}

	img, err := h.imagesService.UploadAndResize(
		ctx,
		authorization,
		data.Name,
		data.Format,
		originalFileHeader,
		croppedFileHeader,
	)
	if err != nil {
		return nil, err
	}

	return http_util.NewResponse(img).WithStatus(http.StatusCreated), nil
}

func (h ImageHandler) deleteOne(ctx context.Context, req *http.Request) (*http_util.Response, error) {
	imageId := chi.URLParam(req, "imageId")
	authDto, err := auth.ExtractAuthorizationDto(ctx, keys.UserAuthDtoKey)
	if err != nil {
		return nil, err
	}
	err = h.imagesService.DeleteOne(ctx, authDto, imageId)
	if err != nil {
		return nil, err
	}

	return http_util.NewResponse(nil).WithStatus(http.StatusNoContent), nil
}
