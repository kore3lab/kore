package helm

import (
	"bytes"

	"helm.sh/helm/v3/pkg/chart"
)

type Renderer struct {
	namespace     string
	componentName string
	chart         *chart.Chart
	buffer        *bytes.Buffer
	started       bool
}

func NewRenderer(componentName, namespace string, bf *bytes.Buffer) *Renderer {

	return &Renderer{
		componentName: componentName,
		namespace:     namespace,
		buffer:        bf,
		started:       false,
	}
}
