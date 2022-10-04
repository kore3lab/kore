package helmreconciler

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"k8s.io/apimachinery/pkg/types"

	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
	"kore3lab.io/kore/pkg/helm"
	"kore3lab.io/kore/pkg/objects"
)

// begin  reconcile
func (o *HelmReconciler) Begin() error {

	if o.operator.Status == nil {
		o.operator.Status = &installv1alpha1.KoreOperatorStatus{Status: installv1alpha1.STATUS_RECONCILING}
	} else {
		o.operator.Status.Status = installv1alpha1.STATUS_RECONCILING
	}

	return nil

}

// end  reconcile
func (o *HelmReconciler) End(status *installv1alpha1.KoreOperatorStatus) error {

	op := &installv1alpha1.KoreOperator{}
	objKey := types.NamespacedName{Name: o.operator.Name, Namespace: o.operator.Namespace}
	if err := o.client.Get(context.TODO(), objKey, op); err != nil {
		return fmt.Errorf("failed to get IstioOperator before updating status due to %v", err)
	}

	if reflect.DeepEqual(o.operator.Status, status) == false {
		o.operator.Status = status
		return o.client.Status().Update(context.TODO(), o.operator)
	}
	return nil
}

// reconcile
func (o *HelmReconciler) Reconcile() *installv1alpha1.KoreOperatorStatus {

	status := &installv1alpha1.KoreOperatorStatus{Status: installv1alpha1.STATUS_RECONCILING}
	if manifest, err := helm.RenderToYaml(o.operator, o.values); err != nil {
		return &installv1alpha1.KoreOperatorStatus{Status: installv1alpha1.STATUS_ERROR, Message: err.Error()}
	} else if k8sObjects, err := objects.GetK8sObjectsFromYaml(manifest); err != nil {
		return &installv1alpha1.KoreOperatorStatus{Status: installv1alpha1.STATUS_ERROR, Message: err.Error()}
	} else {
		var mu sync.Mutex     // mu protects the shared InstallStatus componentStatus across goroutines
		var wg sync.WaitGroup // wg waits for all manifest processing goroutines to finish
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			for _, obj := range *k8sObjects {
				if err := o.ApplyObject(&obj); err != nil {
					status = &installv1alpha1.KoreOperatorStatus{
						Status:  installv1alpha1.STATUS_ERROR,
						Message: fmt.Sprintf("Fail to apply for %s (cause=%v)", obj.GetName(), err),
					}
					break
				}
			}
			mu.Unlock()
		}()

		wg.Wait()
		if status.Status != installv1alpha1.STATUS_ERROR {
			status.Status = installv1alpha1.STATUS_COMPLETE
		}
	}

	return status
}

func (o *HelmReconciler) Finalize() error {

	var mu sync.Mutex     // mu protects the shared InstallStatus componentStatus across goroutines
	var wg sync.WaitGroup // wg waits for all manifest processing goroutines to finish
	for _, name := range []string{"dashboard"} {
		name := name
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()

			if manifest, err := helm.RenderToYaml(o.operator, nil); err != nil {
				o.logger.Info("Fail to render", "name", name, "error", err)
			} else {
				if k8sObjects, err := objects.GetK8sObjectsFromYaml(manifest); err != nil {
					o.logger.Info("Fail to render yaml", "name", name, "error", err)
				} else {
					for _, obj := range *k8sObjects {
						if err := o.DeleteObject(&obj); err != nil {
							o.logger.Info("Fail to delete", "name", name, "error", err)
						}
					}
				}
			}
			mu.Unlock()
		}()

	}
	wg.Wait()
	return nil
}
