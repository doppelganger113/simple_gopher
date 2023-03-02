package image

type Format string

const (
	JpgFormat  Format = "jpg"
	PngFormat  Format = "png"
	WebpFormat Format = "webp"
)

var SupportedFormats = []Format{JpgFormat, PngFormat, WebpFormat}

func (format Format) IsSupported() bool {
	for _, f := range SupportedFormats {
		if f == format {
			return true
		}
	}

	return false
}

func (format Format) ToContentType() ContentType {
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

type ContentType string

const (
	JpegType ContentType = "image/jpeg"
	PngType  ContentType = "image/png"
	WebpType ContentType = "image/webp"
)
