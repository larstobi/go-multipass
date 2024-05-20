package multipass

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindByAlias(t *testing.T) {
	result, err := FindByAlias("jammy")

	if err != nil {
		t.Fatal(err)
	}

	expected := Image{
		Os:      "Ubuntu",
		Release: "22.04 LTS",
		Remote:  "",
		Version: "20220712",
		Aliases: []string{},
	}

	imagesEqual := true

	if result.Os != expected.Os {
		imagesEqual = false
	}
	if result.Release != expected.Release {
		imagesEqual = false
	}
	if reflect.DeepEqual(result.Aliases, expected.Aliases) && imagesEqual {
		t.Logf("Success !")
	} else {
		fmt.Println("Expected:")
		fmt.Println(expected)
		fmt.Println("Got:")
		fmt.Println(*result)
		t.Errorf("Failed")
	}
}
