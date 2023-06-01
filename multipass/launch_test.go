package multipass

import (
	"testing"
)

func TestLaunch(t *testing.T) {

	instanceName := "instanceName"

	defer func() {
		if r := recover(); r != nil {
			// Clean up after testing failure
			Delete(&DeleteRequest{
				Name: instanceName,
			})
			panic(r)
		}
	}()

	instance, err := LaunchV2(&LaunchReqV2{
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
