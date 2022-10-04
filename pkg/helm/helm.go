package helm

import (
	"bytes"

	"helm.sh/helm/v3/pkg/strvals"

	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
)

// get a chart yaml
func RenderChartToYaml(chartName, namespace string, additionalValues []string) (string, error) {

	// values
	values := make(map[string]interface{})
	for _, s := range additionalValues {
		if err := strvals.ParseInto(s, values); err != nil {
			return "", err
		}
	}

	// render a chart
	buffer := &bytes.Buffer{}
	renderer := NewRenderer(chartName, namespace, buffer)
	if err := renderer.Run(); err != nil {
		return buffer.String(), err
	} else if err := renderer.Render(values); err != nil {
		return buffer.String(), err
	}

	return buffer.String(), nil

}

// get a yaml (by profile)
func RenderToYaml(operator *installv1alpha1.KoreOperator, additionalValues []string) (string, error) {

	// merge values
	for _, s := range additionalValues {
		if err := strvals.ParseInto(s, operator.Spec.Values); err != nil {
			return "", err
		}
	}

	// render charts
	buffer := &bytes.Buffer{}
	for name, c := range operator.Spec.Components {
		if c.Enabled == true {
			renderer := NewRenderer(name, operator.GetNamespace(), buffer)
			if err := renderer.Run(); err != nil {
				return buffer.String(), err
			} else if err := renderer.Render(operator.Spec.Values); err != nil {
				return buffer.String(), err
			}
		}
	}

	return buffer.String(), nil

}
