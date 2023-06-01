package multipass

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

const (
	LONG_MEMORY_VERSION = "v1.11.0.rc"
)
