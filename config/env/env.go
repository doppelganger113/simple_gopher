package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	tag         = "env"
	tagName     = "name"
	tagRequired = "required"
	tagDefault  = "default"
)

type tagData struct {
	name         string
	required     bool
	defaultValue string
}

func newTagDataFromTagString(tagString string) (tagData, error) {
	tagMetadata := tagData{}

	if tagString == "" {
		return tagMetadata, nil
	}

	parts := strings.Split(tagString, ",")
	for _, part := range parts {

		keyValueParts := strings.Split(part, ":")
		key := keyValueParts[0]

		switch key {
		case tagName:
			if len(keyValueParts) != 2 {
				return tagData{}, fmt.Errorf("%s tag should have key:value field for '%s'", tag, part)
			}
			tagMetadata.name = keyValueParts[1]
		case tagDefault:
			if len(keyValueParts) != 2 {
				return tagData{}, fmt.Errorf("%s tag should have key:value field for '%s'", tag, part)
			}
			tagMetadata.defaultValue = keyValueParts[1]
		case tagRequired:
			if len(keyValueParts) != 2 {
				tagMetadata.required = true
			} else {
				tagMetadata.required = keyValueParts[1] == "true"
			}
		}
	}

	return tagMetadata, nil
}

func setValue(fieldType reflect.StructField, field reflect.Value, envValue string) error {
	switch fieldType.Type.Kind() {
	case reflect.String:
		field.SetString(envValue)
		return nil
	case reflect.Bool:
		field.SetBool(envValue == "true")
		return nil
	case reflect.Int:
		if envValue == "" {
			return nil
		}
		parsedInt, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(parsedInt)
		return nil
	case reflect.Struct:
		val := field.Addr().Interface()
		errors := UnmarshalFromEnvironment(val)
		if len(errors) > 0 {
			return errors[0]
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %v for field %s", fieldType.Type.Kind(), fieldType.Name)
	}
}

func UnmarshalFromEnvironment(obj interface{}) (errors []error) {
	element := reflect.ValueOf(obj).Elem()
	elementType := reflect.TypeOf(obj).Elem()

	for i := 0; i < element.NumField(); i++ {
		field := element.Field(i)
		fieldType := elementType.Field(i)

		if !fieldType.IsExported() {
			continue
		}
		val, ok := fieldType.Tag.Lookup("env")
		fmt.Printf("env: Value: %s Ok: %v\n", val, ok)

		fmt.Printf("element: %+v \t field: %+v \t fieldType: %+v \t Name: %v\n", element, field, fieldType, fieldType.Name)
		meta, err := newTagDataFromTagString(val)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		fmt.Printf("Metadata %+v\n", meta)
		fmt.Println(fieldType.Name, field.Interface())

		envValue := os.Getenv(meta.name)
		if envValue == "" && meta.defaultValue != "" {
			envValue = meta.defaultValue
		}
		if meta.required && envValue == "" {
			errors = append(errors, fmt.Errorf("missing environment variable %s", meta.name))
			continue
		}
		if setErr := setValue(fieldType, field, envValue); setErr != nil {
			errors = append(errors, setErr)
		}

		fmt.Println("-----------------")
	}

	return
}
