package test

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"os"
	"testing"
)

type Status string

func SkipIfNotIntegrationTesting(t *testing.T) {
	isIntegration := os.Getenv("TEST_INTEGRATION")
	if isIntegration != "true" {
		t.Skip("skipping integration test")
	}
}

type FormParams map[string]string

type FormFilesParams struct {
	FilePath  string
	FileField string
	FileName  string
}

// CreateMultipartFormData creates multipart form data with multiple files and fields
func CreateMultipartFormData(
	formFileParams []FormFilesParams, params FormParams,
) (*bytes.Buffer, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for _, formFileParam := range formFileParams {
		file, err := os.Open(formFileParam.FilePath)
		if err != nil {
			return nil, "", err
		}

		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, "", err
		}
		_ = file.Close()

		part, err := writer.CreateFormFile(formFileParam.FileField, formFileParam.FileName)
		if err != nil {
			return nil, "", err
		}

		_, err = part.Write(fileContents)
		if err != nil {
			return nil, "", err
		}
	}

	for name, value := range params {
		err := writer.WriteField(name, value)
		if err != nil {
			return nil, "", err
		}
	}

	contentType := writer.FormDataContentType()

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, contentType, nil
}
