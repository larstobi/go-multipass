package multipass

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

/*
Public API
*/

type InfoRequest struct {
	Name string
}

func Info(req *InfoRequest) (*Instance, error) {
	if req == nil || req.Name == "" {
		return nil, errors.New("instance name is required")
	}

	cmd := exec.Command(
		"multipass",
		"info",
		req.Name,
		"--format",
		"json",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + " " + err.Error())
	}

	return parseInfoJSON(out, req.Name)
}

/*
Internal JSON parsing
*/

type multipassInfoResponse struct {
	Errors []any                            `json:"errors"`
	Info   map[string]multipassInstanceInfo `json:"info"`
}

type multipassInstanceInfo struct {
	CPUCount  string    `json:"cpu_count"` // i dumpen din er dette string
	State     string    `json:"state"`
	IPv4      []string  `json:"ipv4"`
	Release   string    `json:"release"`
	ImageHash string    `json:"image_hash"`
	Load      []float64 `json:"load"`

	Memory struct {
		Total uint64 `json:"total"`
		Used  uint64 `json:"used"`
	} `json:"memory"`

	Disks map[string]struct {
		Total string `json:"total"` // i dumpen: string med bytes
		Used  string `json:"used"`  // i dumpen: string med bytes
	} `json:"disks"`
}

func parseInfoJSON(data []byte, name string) (*Instance, error) {
	var resp multipassInfoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	info, ok := resp.Info[name]
	if !ok {
		return nil, errors.New("instance not found in multipass output")
	}

	inst := &Instance{
		Name:      name,
		CPUS:      info.CPUCount,
		State:     info.State,
		Image:     info.Release,
		ImageHash: info.ImageHash,
	}

	if len(info.IPv4) > 0 {
		inst.IP = info.IPv4[0]
	}

	if len(info.Load) == 3 {
		inst.Load = fmt.Sprintf("%.2f %.2f %.2f", info.Load[0], info.Load[1], info.Load[2])
	}

	// Disk: velg første disk-entry (typisk sda1). Hvis du vil ha en bestemt disk (root),
	// må du definere policy (f.eks. alltid "sda1" hvis finnes).
	for _, d := range info.Disks {
		inst.DiskUsage = d.Used
		inst.TotalDisk = d.Total
		break
	}

	// Memory i bytes -> lagre som bytes-string (evt. gjør humanize senere)
	inst.MemoryUsage = fmt.Sprintf("%d", info.Memory.Used)
	inst.MemoryTotal = fmt.Sprintf("%d", info.Memory.Total)

	return inst, nil
}
