package multipass

import (
	"encoding/json"
	"errors"
	"os/exec"
)

type FindResult struct {
	Errors []string         `json:"errors"`
	Images map[string]Image `json:"images"`
}

type Image struct {
	Os      string   `json:"os"`
	Release string   `json:"release"`
	Remote  string   `json:"remote"`
	Version string   `json:"version"`
	Aliases []string `json:"aliases"`
}

func FindByAlias(alias string) (*Image, error) {
	args := []string{"find"}
	args = append(args, alias, "--format", "json")

	result := exec.Command("multipass", args...)
	out, err := result.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + " " + err.Error())
	}

	var images FindResult
	json_err := json.Unmarshal(out, &images)
	if json_err != nil {
		return nil, errors.New("Could not unmarshal JSON: " + err.Error())
	}

	var image Image
	image = images.Images[alias]

	return &image, nil
}
