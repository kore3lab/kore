package objects

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Yaml(string) to K8sObject array
func GetK8sObjectsFromYaml(manifest string) (*[]unstructured.Unstructured, error) {

	var buffer bytes.Buffer
	k8sObjects := &[]unstructured.Unstructured{}

	// separate to k8s objects
	scanner := bufio.NewScanner(strings.NewReader(manifest))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") { // yaml separator
			if k8sObject, err := toK8sObject(buffer.Bytes()); err == nil {
				*k8sObjects = append(*k8sObjects, *k8sObject)
			} else {
				return nil, err
			}
			buffer.Reset()
		} else {
			if _, err := buffer.WriteString(line); err != nil {
				return nil, err
			}
			if _, err := buffer.WriteString("\n"); err != nil {
				return nil, err
			}
		}
	}
	if len(buffer.Bytes()) > 0 {
		if k8sObject, err := toK8sObject(buffer.Bytes()); err == nil {
			*k8sObjects = append(*k8sObjects, *k8sObject)
		} else {
			return nil, err
		}
	}

	return k8sObjects, nil

}

func toK8sObject(d []byte) (*unstructured.Unstructured, error) {

	obj := &unstructured.Unstructured{}
	r := bytes.NewReader(d)
	decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)
	for {
		// payload 읽기
		if err := decoder.Decode(obj); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

	}

	return obj, nil
}
