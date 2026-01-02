package multipass

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestParseInfoJSON_OK_AllFields(t *testing.T) {
	t.Parallel()

	const name = "testInfoTest"

	// Matcher ekte format: { "errors": [], "info": { "<name>": { ... } } }
	// Merk: disks.total/used er strings (bytes), memory.total/used er tall (bytes)
	jsonOut := []byte(`{
	  "errors": [],
	  "info": {
		"testInfoTest": {
		  "cpu_count": "2",
		  "state": "Running",
		  "ipv4": ["192.168.2.26"],
		  "release": "Ubuntu 24.04.3 LTS",
		  "image_hash": "a40713938d74",
		  "load": [0.10, 0.05, 0.02],
		  "memory": { "used": 245563392, "total": 3100508160 },
		  "disks": { "sda1": { "used": "2184481280", "total": "5119532032" } }
		}
	  }
	}`)

	inst, err := parseInfoJSON(jsonOut, name)
	if err != nil {
		t.Fatalf("parseInfoJSON error: %v", err)
	}

	if inst.Name != name {
		t.Fatalf("Name: expected %q, got %q", name, inst.Name)
	}
	if inst.State != "Running" {
		t.Fatalf("State: expected %q, got %q", "Running", inst.State)
	}
	if inst.IP != "192.168.2.26" {
		t.Fatalf("IP: expected %q, got %q", "192.168.2.26", inst.IP)
	}
	if inst.Image != "Ubuntu 24.04.3 LTS" {
		t.Fatalf("Image: expected %q, got %q", "Ubuntu 24.04.3 LTS", inst.Image)
	}
	if inst.ImageHash != "a40713938d74" {
		t.Fatalf("ImageHash: expected %q, got %q", "a40713938d74", inst.ImageHash)
	}
	if inst.Load != "0.10 0.05 0.02" {
		t.Fatalf("Load: expected %q, got %q", "0.10 0.05 0.02", inst.Load)
	}

	// DiskUsage/TotalDisk kommer fra første disk-entry (vi har kun én i testen)
	if inst.DiskUsage != "2184481280" || inst.TotalDisk != "5119532032" {
		t.Fatalf("Disk: expected used=%q total=%q, got used=%q total=%q",
			"2184481280", "5119532032", inst.DiskUsage, inst.TotalDisk)
	}

	// MemoryUsage/MemoryTotal er fmt.Sprintf("%d", uint64)
	if inst.MemoryUsage != "245563392" || inst.MemoryTotal != "3100508160" {
		t.Fatalf("Memory: expected used=%q total=%q, got used=%q total=%q",
			"245563392", "3100508160", inst.MemoryUsage, inst.MemoryTotal)
	}
}

func TestParseInfoJSON_OK_MissingIPv4AndNon3Load(t *testing.T) {
	t.Parallel()

	const name = "x"
	jsonOut := []byte(`{
	  "errors": [],
	  "info": {
		"x": {
		  "cpu_count": "1",
		  "state": "Stopped",
		  "ipv4": [],
		  "release": "Ubuntu",
		  "image_hash": "hash",
		  "load": [0.1, 0.2],
		  "memory": { "used": 0, "total": 1 },
		  "disks": { "sda1": { "used": "0", "total": "10" } }
		}
	  }
	}`)

	inst, err := parseInfoJSON(jsonOut, name)
	if err != nil {
		t.Fatalf("parseInfoJSON error: %v", err)
	}

	if inst.IP != "" {
		t.Fatalf("IP: expected empty, got %q", inst.IP)
	}
	// Load skal ikke settes hvis len(load) != 3
	if inst.Load != "" {
		t.Fatalf("Load: expected empty, got %q", inst.Load)
	}
}

func TestParseInfoJSON_Err_InstanceNotFound(t *testing.T) {
	t.Parallel()

	jsonOut := []byte(`{
	  "errors": [],
	  "info": {
		"someOther": { "state": "Running" }
	  }
	}`)

	_, err := parseInfoJSON(jsonOut, "missing")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "instance not found") {
		t.Fatalf("expected instance not found error, got: %v", err)
	}
}

func TestParseInfoJSON_Err_InvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := parseInfoJSON([]byte(`{not-json`), "x")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestInfo_UsesMultipassJSON_OK(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("this test uses a temporary shell script multipass; skip on Windows")
	}

	tmp, err := ioutil.TempDir("", "go-multipass-test-")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tmp)

	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Fake multipass: returnerer *ekte* JSON-format (errors/info/<name>)
	script := `#!/bin/sh
# Expect: multipass info <name> --format json
if [ "$1" != "info" ]; then
  echo "unexpected args" >&2
  exit 2
fi

NAME="$2"

cat <<EOF
{
  "errors": [],
  "info": {
    "$NAME": {
      "cpu_count": "2",
      "state": "Running",
      "ipv4": ["192.168.2.26"],
      "release": "Ubuntu 24.04.3 LTS",
      "image_hash": "a40713938d74",
      "load": [0.10, 0.05, 0.02],
      "memory": { "used": 245563392, "total": 3100508160 },
      "disks": { "sda1": { "used": "2184481280", "total": "5119532032" } }
    }
  }
}
EOF
`

	multipassPath := filepath.Join(binDir, "multipass")
	if err := ioutil.WriteFile(multipassPath, []byte(script), 0755); err != nil {
		t.Fatalf("write fake multipass: %v", err)
	}

	// Prepend fake multipass to PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	if err := os.Setenv("PATH", binDir+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatalf("set PATH: %v", err)
	}

	inst, err := Info(&InfoRequest{Name: "testInfoTest"})
	if err != nil {
		// Ekstra diagnose: kjør fake multipass direkte for å se hva Info() faktisk får
		out, _ := exec.Command("multipass", "info", "testInfoTest", "--format", "json").CombinedOutput()
		t.Fatalf("Info error: %v\nfake multipass output:\n%s", err, string(out))
	}

	if inst.Name != "testInfoTest" {
		t.Fatalf("Name: expected %q, got %q", "testInfoTest", inst.Name)
	}
	if inst.IP != "192.168.2.26" {
		t.Fatalf("IP: expected %q, got %q", "192.168.2.26", inst.IP)
	}
	if inst.Load != "0.10 0.05 0.02" {
		t.Fatalf("Load: expected %q, got %q", "0.10 0.05 0.02", inst.Load)
	}
}

func TestInfo_PropagatesCommandError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("this test uses a temporary shell script multipass; skip on Windows")
	}

	tmp, err := ioutil.TempDir("", "go-multipass-test-")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tmp)

	binDir := filepath.Join(tmp, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	script := `#!/bin/sh
echo "boom from multipass" >&2
exit 42
`
	multipassPath := filepath.Join(binDir, "multipass")
	if err := ioutil.WriteFile(multipassPath, []byte(script), 0755); err != nil {
		t.Fatalf("write fake multipass: %v", err)
	}

	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	if err := os.Setenv("PATH", binDir+string(os.PathListSeparator)+oldPath); err != nil {
		t.Fatalf("set PATH: %v", err)
	}

	_, err = Info(&InfoRequest{Name: "any"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "boom from multipass") {
		t.Fatalf("expected stderr to be included, got: %v", err)
	}
}
