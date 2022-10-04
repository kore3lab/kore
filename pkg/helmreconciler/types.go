package helmreconciler

import (
	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
)

type HelmReconciler struct {
	operator *installv1alpha1.KoreOperator
	values   []string
	client   client.Client
	logger   logr.Logger
}

func NewHelmReconciler(op *installv1alpha1.KoreOperator, additionalValues []string, c client.Client) (*HelmReconciler, error) {

	if c == nil {
		scheme := runtime.NewScheme()
		utilruntime.Must(clientgoscheme.AddToScheme(scheme))
		utilruntime.Must(installv1alpha1.AddToScheme(scheme))
		if mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{Scheme: scheme}); err != nil {
			return nil, err
		} else {
			c = mgr.GetClient()
		}
	}

	return (&HelmReconciler{
		operator: op,
		client:   c,
		logger:   ctrl.Log.WithName("scoped"),
		values:   additionalValues,
	}), nil
}
