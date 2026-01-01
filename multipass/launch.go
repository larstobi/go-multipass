package multipass

import (
	"errors"
	"os/exec"
	"strings"
)

type LaunchReq struct {
	Image         string
	CPUS          string
	Disk          string
	Name          string
	Memory        string
	CloudInitFile string
	CloudInitData string
	CloudInitUrl  string
	Network       []string
	Bridged       bool
}

func Launch(launchReq *LaunchReq) (*Instance, error) {
	if launchReq == nil {
		return nil, errors.New("launchReq is nil")
	}

	// Enforce mutual exclusivity
	set := 0
	if launchReq.CloudInitFile != "" {
		set++
	}
	if launchReq.CloudInitData != "" {
		set++
	}
	if launchReq.CloudInitUrl != "" {
		set++
	}
	if set > 1 {
		return nil, errors.New(
			"only one of CloudInitFile, CloudInitData, or CloudInitUrl can be set",
		)
	}

	args := []string{"launch"}

	if launchReq.Image != "" {
		args = append(args, launchReq.Image)
	}

	if launchReq.CPUS != "" {
		args = append(args, "--cpus", launchReq.CPUS)
	}

	if launchReq.Name != "" {
		args = append(args, "--name", launchReq.Name)
	}

	if launchReq.Disk != "" {
		args = append(args, "--disk", launchReq.Disk)
	}

	if launchReq.Memory != "" {
		args = append(args, "-m", launchReq.Memory)
	}

	// Cloud-init handling
	var stdin *strings.Reader
	var useStdin bool

	switch {
	case launchReq.CloudInitFile != "":
		args = append(args, "--cloud-init", launchReq.CloudInitFile)

	case launchReq.CloudInitUrl != "":
		args = append(args, "--cloud-init", launchReq.CloudInitUrl)

	case launchReq.CloudInitData != "":
		args = append(args, "--cloud-init", "-")

		data := launchReq.CloudInitData
		if !strings.HasSuffix(data, "\n") {
			data += "\n"
		}
		stdin = strings.NewReader(data)
		useStdin = true
	}

	// Network specs
	for _, network := range launchReq.Network {
		args = append(args, "--network", network)
	}
	if launchReq.Bridged {
		args = append(args, "--bridged")
	}

	cmd := exec.Command("multipass", args...)
	if useStdin {
		cmd.Stdin = stdin
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + " " + err.Error())
	}

	out2 := strings.TrimSpace(string(out))
	if out2 == "" {
		return nil, errors.New("empty multipass output")
	}
	lines := strings.Split(out2, "\n")

	// Expect: "Launched: <name>"
	const launchedPrefix = "Launched: "
	line := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(line, launchedPrefix) {
		return nil, errors.New("unexpected multipass output: " + out2)
	}
	name := strings.TrimSpace(strings.TrimPrefix(line, launchedPrefix))

	instance, err := Info(&InfoRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return instance, nil
}
