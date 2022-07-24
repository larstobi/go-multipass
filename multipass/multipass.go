package multipass

import "strings"

const (
	Name        = "Name:"
	State       = "State:"
	IPv4        = "IPv4:"
	Release     = "Release:"
	ImageHash   = "Image hash:"
	Load        = "Load:"
	DiskUsage   = "Disk usage:"
	MemoryUsage = "Memory usage:"
)

type Instance struct {
	Name        string
	IP          string
	State       string
	Image       string
	ImageHash   string
	DiskUsage   string
	TotalDisk   string
	MemoryUsage string
	MemoryTotal string
	Load        string
}

func parseInstance(out string) *Instance {

	var instance Instance

	for _, line := range strings.Split(out, "\n") {

		if strings.Contains(line, Name) && !strings.HasSuffix(line, "--") {
			instance.Name = strings.TrimSpace(strings.ReplaceAll(line, Name, ""))
		}

		if strings.Contains(line, State) && !strings.HasSuffix(line, "--") {
			instance.State = strings.TrimSpace(strings.ReplaceAll(line, State, ""))
		}

		if strings.Contains(line, IPv4) && !strings.HasSuffix(line, "--") {
			instance.IP = strings.TrimSpace(strings.ReplaceAll(line, IPv4, ""))
		}

		if strings.Contains(line, Release) && !strings.HasSuffix(line, "--") {
			instance.Image = strings.TrimSpace(strings.ReplaceAll(line, Release, ""))
		}

		if strings.Contains(line, ImageHash) && !strings.HasSuffix(line, "--") {
			instance.ImageHash = strings.TrimSpace(strings.ReplaceAll(line, ImageHash, ""))
		}

		if strings.Contains(line, Load) && !strings.HasSuffix(line, "--") {
			instance.Load = strings.TrimSpace(strings.ReplaceAll(line, Load, ""))
		}

		if strings.Contains(line, DiskUsage) && !strings.HasSuffix(line, "--") {
			diskUsage := strings.TrimSpace(strings.ReplaceAll(line, DiskUsage, ""))
			diskUsageOut := strings.Split(diskUsage, "out of")
			instance.DiskUsage = strings.TrimSpace(diskUsageOut[0])
			instance.TotalDisk = strings.TrimSpace(diskUsageOut[1])
		}

		if strings.Contains(line, MemoryUsage) && !strings.HasSuffix(line, "--") {
			memoryUsage := strings.TrimSpace(strings.ReplaceAll(line, MemoryUsage, ""))
			memoryUsageOut := strings.Split(memoryUsage, "out of")
			instance.MemoryUsage = strings.TrimSpace(memoryUsageOut[0])
			instance.MemoryTotal = strings.TrimSpace(memoryUsageOut[1])
		}
	}

	return &instance
}
