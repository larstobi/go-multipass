package multipass

import (
	"encoding/json"
	"os/exec"
)

type multipassListResponse struct {
	List []multipassInstance `json:"list"`
}

type multipassInstance struct {
	Name  string `json:"name"`
	State string `json:"state"`
	// Behold gjerne disse i JSON-strukturen, selv om Instance ikke har dem:
	IPv4    []string `json:"ipv4"`
	Release string   `json:"release"`
}

func List() ([]*Instance, error) {
	cmd := exec.Command("multipass", "list", "--format", "json")

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var resp multipassListResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil, err
	}

	instances := make([]*Instance, 0, len(resp.List))
	for _, mp := range resp.List {
		instances = append(instances, &Instance{
			Name:  mp.Name,
			State: mp.State,
		})
	}

	return instances, nil
}
