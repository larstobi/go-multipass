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

	if reflect.DeepEqual(*result, expected) {
		t.Logf("Success !")
	} else {
		fmt.Println("Expected:")
		fmt.Println(expected)
		fmt.Println("Got:")
		fmt.Println(*result)
		t.Errorf("Failed")
	}
}
