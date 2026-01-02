package multipass

import (
	"os/exec"
	"testing"
)

func dumpMultipassState(t *testing.T, name string) {
	t.Helper()

	// Best-effort: ingen fail i dump-funksjonen
	if out, err := exec.Command("multipass", "info", name, "--format", "json").CombinedOutput(); err == nil {
		t.Logf("DEBUG: multipass info %q:\n%s", name, string(out))
	} else {
		t.Logf("DEBUG: multipass info %q failed: %v\n%s", name, err, string(out))
	}

	if out, err := exec.Command("multipass", "list", "--format", "json").CombinedOutput(); err == nil {
		t.Logf("DEBUG: multipass list:\n%s", string(out))
	} else {
		t.Logf("DEBUG: multipass list failed: %v\n%s", err, string(out))
	}

	if out, err := exec.Command("multipass", "get", "local.driver").CombinedOutput(); err == nil {
		t.Logf("DEBUG: multipass get local.driver:\n%s", string(out))
	}
}

func TestLaunch_NilRequest(t *testing.T) {
	instance, err := Launch(nil)

	if err == nil {
		t.Fatalf("expected error when launchReq is nil, got nil")
	}

	if instance != nil {
		t.Fatalf("expected nil instance when launchReq is nil, got %#v", instance)
	}
}

func TestLaunch(t *testing.T) {

	instanceName := "testTestLaunch"

	defer func() {
		if t.Failed() {
			dumpMultipassState(t, instanceName)
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			// Clean up after testing failure
			Delete(&DeleteRequest{
				Name: instanceName,
			})
			panic(r)
		}
	}()

	instance, err := Launch(&LaunchReq{
		CPUS:   "2",
		Memory: "3G",
		Name:   instanceName,
	})
	if err != nil {
		t.Fatal(err)
	} else {
		if instance.MemoryTotal != "2.9GiB" {
			t.Error("Expected memory setting: 2.9GiB, got: " + instance.MemoryTotal)
		}
	}

	// Clean up after testing
	Delete(&DeleteRequest{
		Name: instanceName,
	})
}

func TestLaunch_CloudInitMutualExclusivity_FileAndData(t *testing.T) {
	_, err := Launch(&LaunchReq{
		CloudInitFile: "cloud-init.yaml",
		CloudInitData: "#cloud-config\n",
	})
	if err == nil {
		t.Fatalf("expected error when both CloudInitFile and CloudInitData are set")
	}
}

func TestLaunch_CloudInitMutualExclusivity_FileAndUrl(t *testing.T) {
	_, err := Launch(&LaunchReq{
		CloudInitFile: "cloud-init.yaml",
		CloudInitURL:  "https://example.com/cloud-init.yaml",
	})
	if err == nil {
		t.Fatalf("expected error when both CloudInitFile and CloudInitURL are set")
	}
}

func TestLaunch_CloudInitMutualExclusivity_DataAndUrl(t *testing.T) {
	_, err := Launch(&LaunchReq{
		CloudInitData: "#cloud-config\n",
		CloudInitURL:  "https://example.com/cloud-init.yaml",
	})
	if err == nil {
		t.Fatalf("expected error when both CloudInitData and CloudInitURL are set")
	}
}

func TestLaunch_CloudInitMutualExclusivity_AllThree(t *testing.T) {
	_, err := Launch(&LaunchReq{
		CloudInitFile: "cloud-init.yaml",
		CloudInitData: "#cloud-config\n",
		CloudInitURL:  "https://example.com/cloud-init.yaml",
	})
	if err == nil {
		t.Fatalf("expected error when CloudInitFile, CloudInitData, and CloudInitURL are all set")
	}
}
