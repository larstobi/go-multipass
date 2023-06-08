package multipass

import (
	"errors"
	"os/exec"
	"strings"
)

type LaunchReq struct {
	Image         string
	CPU           string
	Disk          string
	Name          string
	Memory        string
	CloudInitFile string
}

type LaunchReqV2 struct {
	Image         string
	CPUS          string
	Disk          string
	Name          string
	Memory        string
	CloudInitFile string
}

func Launch(launchReq *LaunchReq) (*Instance, error) {
	instance, err := LaunchV2(&LaunchReqV2{
		Name:          launchReq.Name,
		Image:         launchReq.Image,
		CPUS:          launchReq.CPU,
		Memory:        launchReq.Memory,
		Disk:          launchReq.Disk,
		CloudInitFile: launchReq.CloudInitFile,
	})
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func LaunchV2(launchReqV2 *LaunchReqV2) (*Instance, error) {

	var args = []string{"launch"}

	if launchReqV2.Image != "" {
		args = append(args, launchReqV2.Image)
	}

	if launchReqV2.CPUS != "" {
		args = append(args, "--cpus", launchReqV2.CPUS)
	}

	if launchReqV2.Name != "" {
		args = append(args, "--name", launchReqV2.Name)
	}

	if launchReqV2.Disk != "" {
		args = append(args, "--disk", launchReqV2.Disk)
	}

	if launchReqV2.Memory != "" {
		args = append(args, "-m", launchReqV2.Memory)
	}

	if launchReqV2.CloudInitFile != "" {
		args = append(args, "--cloud-init", launchReqV2.CloudInitFile)
	}

	result := exec.Command("multipass", args...)
	out, err := result.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + " " + err.Error())
	}

	var b []byte
	b = append(b, out...)

	out2 := string(b)

	o := strings.Split(strings.TrimSpace(out2), "\n")[0]

	name := strings.TrimSpace(strings.Split(o, "Launched: ")[1])

	instance, err := Info(&InfoRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return instance, nil
}
