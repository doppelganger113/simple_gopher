package http_transport

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"simple_gopher"
	"simple_gopher/auth"
	"simple_gopher/image_resize"
	"simple_gopher/storage"
)

const maxBodyLimitBytes = 30 * 1024 * 1024 // 20MB

func ImagesRouter(handler ImagesHandler, authenticator Authenticator) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", FetchImages(handler))
		r.Get("/{imageId}", FetchImage(handler))
		r.Post("/",
			Authorize(AddImage(handler), authenticator, auth.RoleAdmin),
		)
		r.Patch("/{imageId}",
			Authorize(UpdateImage(handler), authenticator, auth.RoleAdmin),
		)
		r.Delete("/{imageId}",
			Authorize(DeleteOne(handler), authenticator, auth.RoleAdmin),
		)
	}
}

func FetchImage(handler ImagesHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		imageId := chi.URLParam(r, "imageId")
		image, err := handler.GetOne(ctx, imageId)
		if err != nil {
			handleError(w, err)
			return
		}

		respondJson(w, http.StatusOK, image)
	}
}

func FetchImages(handler ImagesHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		page := ToUint(r.URL.Query().Get("page"))
		size := ToUint(r.URL.Query().Get("size"))

		order := storage.ToOrderOr(r.URL.Query().Get("order"), storage.OrderDescending)
		limit, offset := storage.PagingToLimitOffset(page, size)

		imageList, err := handler.Get(ctx, limit, offset, order)
		if err != nil {
			handleError(w, err)
			return
		}

		respondJson(w, http.StatusOK, imageList)
	}
}

type UploadImageDto struct {
	Name   string
	Format image_resize.ImageFormat
}

func (dto UploadImageDto) validate() error {
	if len(dto.Name) < 5 || len(dto.Name) > 200 {
		return simple_gopher.InvalidArgument{
			Reason: "Name should be between 5 and 250 characters",
		}
	}

	if !dto.Format.IsSupported() {
		return simple_gopher.InvalidArgument{
			Reason: fmt.Sprintf("Unsupported format %s", dto.Format),
		}
	}

	return nil
}

func AddImage(handler ImagesHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(maxBodyLimitBytes)

		_, originalFileHeader, err := r.FormFile("originalFile")
		if err != nil {
			respondJson(
				w,
				http.StatusBadRequest,
				newFailureResponse("missing originalFile"),
			)
			return
		}
		_, croppedFileHeader, err := r.FormFile("croppedFile")
		if err != nil {
			respondJson(
				w,
				http.StatusBadRequest,
				newFailureResponse("missing croppedFile"),
			)
			return
		}

		data := &UploadImageDto{}
		data.Name = r.PostFormValue("name")
		data.Format = image_resize.ImageFormat(r.PostFormValue("format"))
		if err = data.validate(); err != nil {
			respondBadRequestJson(w, err)
			return
		}

		ctx := r.Context()
		authorization, err := auth.ExtractAuthorizationDto(ctx, UserAuthDtoKey)
		if err != nil {
			handleError(w, err)
			return
		}

		image, err := handler.UploadAndResize(
			ctx,
			authorization,
			data.Name,
			data.Format,
			originalFileHeader,
			croppedFileHeader,
		)
		if err != nil {
			handleError(w, err)
			return
		}

		respondJson(w, http.StatusCreated, image)
	}
}

func UpdateImage(handler ImagesHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(maxBodyLimitBytes)

		_, originalFileHeader, err := r.FormFile("originalFile")
		if err != nil {
			respondJson(
				w,
				http.StatusBadRequest,
				newFailureResponse("missing originalFile"),
			)
			return
		}
		_, croppedFileHeader, err := r.FormFile("croppedFile")
		if err != nil {
			respondJson(
				w,
				http.StatusBadRequest,
				newFailureResponse("missing croppedFile"),
			)
			return
		}

		data := &UploadImageDto{}
		data.Name = r.PostFormValue("name")
		data.Format = image_resize.ImageFormat(r.PostFormValue("format"))
		if err = data.validate(); err != nil {
			respondBadRequestJson(w, err)
			return
		}

		ctx := r.Context()
		authorization, err := auth.ExtractAuthorizationDto(ctx, UserAuthDtoKey)
		if err != nil {
			handleError(w, err)
			return
		}

		image, err := handler.UploadAndResize(
			ctx,
			authorization,
			data.Name,
			data.Format,
			originalFileHeader,
			croppedFileHeader,
		)
		if err != nil {
			handleError(w, err)
			return
		}

		respondJson(w, http.StatusCreated, image)
	}
}

func DeleteOne(handler ImagesHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		imageId := chi.URLParam(r, "imageId")
		authDto, err := auth.ExtractAuthorizationDto(ctx, UserAuthDtoKey)
		if err != nil {
			handleError(w, err)
			return
		}
		err = handler.DeleteOne(ctx, authDto, imageId)
		if err != nil {
			handleError(w, err)
			return
		}

		respondJson(w, http.StatusNoContent, nil)
	}
}
