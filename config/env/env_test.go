package env

import (
	"fmt"
	"testing"
)

func TestSetFromEnvironment_ShouldValidateCorrectStruct(t *testing.T) {
	type Options struct {
		Pretty       bool   `env:"name:LOG_PRETTY,default:true"`
		Seconds      int    `env:"name:SECONDS,default:20"`
		privateField bool   `env:"name:LOG_PRETTY,default:tester,required"`
		Name         string `env:"name:NAME,required,default:John"`
		Nested       struct {
			Username string `env:"name:USERNAME,required,default:John"`
		}
	}

	options := Options{}
	errors := UnmarshalFromEnvironment(&options)
	if errors != nil {
		for _, err := range errors {
			t.Errorf("%s\n", err)
		}
		t.Fatalf("failed setting from environment, got %d errors", len(errors))
	}

	fmt.Printf("%+v\n", options)
	if !options.Pretty {
		t.Fatal("pretty should be true")
	}
	if options.privateField == true {
		t.Fatal("privateField should not be changed")
	}
	if options.Seconds != 20 {
		t.Fatal("Seconds should be 20")
	}
	if options.Name != "John" {
		t.Fatalf("Name should be John, got: %s", options.Name)
	}
	if options.Nested.Username != "John" {
		t.Fatalf("Nested username should be John, got: %s", options.Nested.Username)
	}
}
