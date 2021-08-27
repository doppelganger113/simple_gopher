package image_resize

type SignedRequest struct {
	Format ImageFormat `json:"format"`
}

type SignedResponse struct {
	SignedUrl string `json:"signedUrl"`
	FileName  string `json:"fileName"`
}

type ImageResizeRequest struct {
	Name             string `json:"name"`
	FilePath         string `json:"filePath"`
	OriginalFilePath string `json:"originalFilePath"`
}

type ImageResizeResponse struct {
	Format   ImageFormat `json:"format"`
	Original string      `json:"original"`
	Name     string      `json:"name"`
	Domain   string      `json:"domain"`
	Path     string      `json:"path"`
	Sizes    ImageSizes  `json:"sizes"`
}

type ImageRenameRequest struct {
	Name    string      `json:"name"`
	NewName string      `json:"newName"`
	Format  ImageFormat `json:"format"`
	SizeMap ImageSizes  `json:"sizeMap"`
}

type ImageDeleteRequest struct {
	Name       string       `json:"name"`
	Format     ImageFormat  `json:"format"`
	Dimensions []Dimensions `json:"dimensions"`
}
