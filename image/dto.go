package image

type SignedRequest struct {
	Format Format `json:"format"`
}

type SignedResponse struct {
	SignedUrl string `json:"signedUrl"`
	FileName  string `json:"fileName"`
}

type ResizeRequest struct {
	Name             string `json:"name"`
	FilePath         string `json:"filePath"`
	OriginalFilePath string `json:"originalFilePath"`
}

type ResizeResponse struct {
	Format   Format `json:"format"`
	Original string `json:"original"`
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Sizes    Sizes  `json:"sizes"`
}

type RenameRequest struct {
	Name    string `json:"name"`
	NewName string `json:"newName"`
	Format  Format `json:"format"`
	SizeMap Sizes  `json:"sizeMap"`
}

type DeleteRequest struct {
	Name       string       `json:"name"`
	Format     Format       `json:"format"`
	Dimensions []Dimensions `json:"dimensions"`
}
