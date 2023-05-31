package multipass

import (
	"log"
	"os/exec"
	"strings"
	"regexp"
	"golang.org/x/mod/semver"
)

func GetVersion() (string, error) {

	cmdString := "multipass --version"

	cmd := exec.Command("sh", "-c", cmdString)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		return "", err
	}
	version := string(out)
	version = regexp.MustCompile("\r?\n").Split(version, -1)[0]
	version = strings.TrimSpace(strings.TrimPrefix(version, "multipass"))
	
	return version, nil
}


func normalizeVersion(version string) string {
	
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	if strings.HasPrefix(version, "v201") { // all versions equals of before 2018.** should be considered older than 0.5
		version = strings.TrimPrefix(version, "v201")
		version = "v0.1.201" + version
	}
	
	//fmt.Println("Normalized version: ", version)
	return version
}

func CompareVersion(v1 string, v2 string) (int) {

	// Compare versions
	normalizedV1 := normalizeVersion(v1)

	normalizedV2 := normalizeVersion(v2)

	result := semver.Compare(normalizedV1, normalizedV2)

	// if result > 0 {
	// 	fmt.Printf("%s is greater than %s\n", v1, v2)
	// } else if result < 0 {
	// 	fmt.Printf("%s is less than %s\n", v1, v2)
	// } else {
	// 	fmt.Printf("%s is equal to %s\n", v1, v2)
	// }

	return result;

}

func CompareCurrentVersion(v2 string) (int) {

	// Compare versions
	currentversion, err := GetVersion()
	if err != nil {
		currentversion = "0.05"
	}
	result := CompareVersion(currentversion, v2)

	return result;

}