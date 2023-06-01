package multipass

import (
	"testing"
)

func TestLaunch(t *testing.T) {

	_, err := LaunchV2(&LaunchReqV2{
		CPUS: "2",
		Memory: "2GiB",
		Name: "instanceName", 
	})
	if err != nil {
		t.Fatal(err)
	}

}
