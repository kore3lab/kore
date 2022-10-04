package manifests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
)

var ManifestsPath = getManifestsPath()

func GetProfile(name string) (*installv1alpha1.KoreOperator, error) {

	// load profile
	operator := &installv1alpha1.KoreOperator{}
	if by, err := ReadFile(name); err != nil {
		return operator, err
	} else {
		if err = yaml.Unmarshal(by, operator); err != nil {
			return operator, err
		}
	}

	return operator, nil

}

// GetProfileYaml reads the YAML values associated with the given profile. It uses an appropriate reader for the
// profile format (compiled-in, file, HTTP, etc.).
func ReadFile(fileName string) ([]byte, error) {

	// Get global values from profile.
	if strings.Contains(fileName, "/") || strings.Contains(fileName, ".") {

		//is FilePath
		if by, err := ioutil.ReadFile(filepath.Join(ManifestsPath, fileName)); err != nil {
			return nil, fmt.Errorf("failed to read profile %v : %v", fileName, err)
		} else {
			return by, nil
		}
	} else {
		// is ProfileName
		if fileName == "" {
			fileName = "default"
		}

		by, err := ioutil.ReadFile(filepath.Join(ManifestsPath, "profiles", fileName+".yaml"))
		if err != nil {
			return nil, fmt.Errorf("failed to read profile %v : %v", fileName, err)
		} else {
			return by, nil
		}
	}

}

// list all the profiles.
func ListProfiles() ([]string, error) {

	// read profiles
	profiles := map[string]bool{}

	files, _ := ioutil.ReadDir(filepath.Join(ManifestsPath, "profiles"))

	for _, f := range files {
		if f.IsDir() == false {
			trimmedString := strings.TrimSuffix(f.Name(), ".yaml")
			if f.Name() != trimmedString {
				profiles[trimmedString] = true
			}
		}
	}

	s := make([]string, 0, len(profiles))
	for k, v := range profiles {
		if v {
			s = append(s, k)
		}
	}

	return s, nil
}

func getManifestsPath() string {

	findDirs := []string{}

	//execute dir
	ex, _ := os.Executable()
	exeDir, _ := filepath.Split(ex)
	if strings.HasSuffix(exeDir, "/bin/") {
		findDirs = append(findDirs, filepath.Join(exeDir, "../"))
	} else {
		findDirs = append(findDirs, exeDir)
	}

	// env. BASE_DIR
	if os.Getenv("BASE_DIR") != "" {
		findDirs = append(findDirs, os.Getenv("BASE_DIR"))
	}

	// pwd
	if dir, err := os.Getwd(); err == nil {
		findDirs = append(findDirs, dir)
	}

	//find a dirctory 'manifests'
	var manifestsPath string
	for _, rootDir := range findDirs {
		if findSubDir(rootDir, "manifests") == true {
			manifestsPath = filepath.Join(rootDir, "manifests")
			break
		}
	}

	return manifestsPath

}

func findSubDir(rootDir, dirName string) bool {

	dirs, _ := ioutil.ReadDir(rootDir)
	for _, dir := range dirs {
		if dir.IsDir() && dir.Name() == dirName {
			return true
		}
	}

	return false
}
