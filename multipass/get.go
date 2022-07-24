package multipass

import (
    "errors"
    "os/exec"
    "strings"
)

type GetReq struct {
    Name string
}

type InstanceAttributes struct {
    CPUS   string
    Disk   string
    Name   string
    Memory string
}

func GetCmd(name string, attribute string) (string, error) {
    args := []string{"get"}
    local := []string{"local", name, attribute}
    args = append(args, strings.Join(local, "."))

    result := exec.Command("multipass", args...)
    out, err := result.CombinedOutput()
    if err != nil {
        return "", errors.New(string(out) + " " + err.Error())
    }

    // Remove trailing newlines from CLI output
    ret := strings.Replace(string(out), "\n", "", -1)

    return ret, nil
}

func Get(getReq *GetReq) (*InstanceAttributes, error) {
    instance := InstanceAttributes{}

    attrs := []string{"cpus", "memory", "disk"}
    var attrs_result = make(map[string]string)

    for _, attr := range attrs {
        out, err := GetCmd(getReq.Name, attr)
        if err != nil {
            return nil, errors.New(string(out) + " " + err.Error())
        }
        attrs_result[attr] = out
    }

    instance.CPUS = attrs_result["cpus"]
    instance.Memory = attrs_result["memory"]
    instance.Disk = attrs_result["disk"]

    return &instance, nil
}
