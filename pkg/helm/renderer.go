package helm

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/apimachinery/pkg/version"

	"kore3lab.io/kore/manifests"
)

const (
	YAMLSeparator = "\n---\n"
)

var KubernetesVersion *version.Info // Kubernetes cluster version

// Run load chart
func (o *Renderer) Run() error {

	var err error
	o.chart, err = loader.LoadDir(filepath.Join(manifests.ManifestsPath, "charts", o.componentName))
	if err != nil {
		return fmt.Errorf("load files: %v", err)
	}
	o.started = true

	return nil
}

// RenderManifest renders the current helm templates with the current values and returns the resulting YAML manifest string.
// - See https://github.com/helm/helm/blob/main/pkg/lint/rules/template.go
func (o *Renderer) Render(values map[string]interface{}) error {

	if !o.started {
		return fmt.Errorf("fileTemplateRenderer for %s not started in renderChart", o.componentName)
	}

	// values
	options := chartutil.ReleaseOptions{Name: o.componentName, Namespace: o.namespace}
	caps := *chartutil.DefaultCapabilities
	if KubernetesVersion != nil {
		caps.KubeVersion = chartutil.KubeVersion{Version: KubernetesVersion.GitVersion, Major: KubernetesVersion.Major, Minor: KubernetesVersion.Minor}
	}

	// subhchart dependencies
	if err := chartutil.ProcessDependencies(o.chart, values); err != nil {
		return err
	}

	vals, err := chartutil.ToRenderValues(o.chart, values, options, &caps)
	if err != nil {
		return err
	}

	// rendering
	outputs, err := engine.Render(o.chart, vals)
	if err != nil {
		return err
	}

	// add yaml separator if the rendered file doesn't have one at the end
	keys := make([]string, 0, len(outputs))
	for k := range outputs {
		if strings.HasSuffix(k, ".txt") {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i < len(keys); i++ {

		name := keys[i]
		data := strings.TrimSpace(outputs[name])
		if data != "" {
			data = "# Source: " + name + "\n" + data
			if !strings.HasSuffix(data, YAMLSeparator) {
				data += YAMLSeparator
			}
			if _, err := o.buffer.WriteString(data); err != nil {
				return err
			}
		}
	}

	return nil
}
