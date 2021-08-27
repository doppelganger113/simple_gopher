package storage

type ImageFormat string

const (
	JpgFormat  ImageFormat = "jpg"
	PngFormat  ImageFormat = "png"
	WebpFormat ImageFormat = "webp"
)

var SupportedImageFormats = []ImageFormat{JpgFormat, PngFormat, WebpFormat}

func (format ImageFormat) IsSupported() bool {
	for _, f := range SupportedImageFormats {
		if f == format {
			return true
		}
	}

	return false
}

func (format ImageFormat) ToImageContentType() ImageContentType {
	switch format {
	case JpgFormat:
		return JpegType
	case PngFormat:
		return PngType
	case WebpFormat:
		return WebpType
	}

	return ""
}

type ImageContentType string

const (
	JpegType ImageContentType = "image/jpeg"
	PngType  ImageContentType = "image/png"
	WebpType ImageContentType = "image/webp"
)
