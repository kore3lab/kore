package helmreconciler

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApplyObject creates or updates an object in the API server depending on whether it already exists.
func (o *HelmReconciler) ApplyObject(obj *unstructured.Unstructured) error {

	if obj.GetKind() == "List" {
		errs := []error{}
		list, err := obj.ToList()
		if err != nil {
			return err
		}
		for _, item := range list.Items {
			err = o.ApplyObject(&item)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return fmt.Errorf("%s", errs)
	}

	objectKey := client.ObjectKeyFromObject(obj)

	receiver := &unstructured.Unstructured{}
	receiver.SetGroupVersionKind(obj.GroupVersionKind())
	err := o.client.Get(context.TODO(), objectKey, receiver)
	if err != nil {
		if apierrors.IsNotFound(err) {
			err = o.client.Create(context.TODO(), obj)
			if err != nil {
				return fmt.Errorf("Failed to create %s/%s/%s: %w", obj.GetKind(), obj.GetNamespace(), obj.GetName(), err)
			} else {
				o.logger.Info("Create a object", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
			}
		} else {
			o.logger.Error(err, "Skip to create a object", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
		}
	} else {
		if err := o.client.Update(context.TODO(), receiver); err != nil {
			return fmt.Errorf("Failed to update %s/%s/%s: %w", obj.GetKind(), obj.GetNamespace(), obj.GetName(), err)
		} else {
			o.logger.Info("Update a object", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
		}
	}
	return nil
}

func (o *HelmReconciler) DeleteObject(obj *unstructured.Unstructured) error {

	if obj.GetKind() == "List" {
		errs := []error{}
		list, err := obj.ToList()
		if err != nil {
			return err
		}
		for _, item := range list.Items {
			err = o.DeleteObject(&item)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return fmt.Errorf("%s", errs)
	}

	objectKey := client.ObjectKeyFromObject(obj)

	receiver := &unstructured.Unstructured{}
	receiver.SetGroupVersionKind(obj.GroupVersionKind())
	if err := o.client.Get(context.TODO(), objectKey, receiver); err == nil {
		err = o.client.Delete(context.TODO(), obj)
		if err != nil {
			return fmt.Errorf("failed to delete %s/%s/%s: %w", obj.GetKind(), obj.GetNamespace(), obj.GetName(), err)
		} else {
			o.logger.Info("Delete a object", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
		}
	} else {
		o.logger.Error(err, "Skip to delete a object", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
	}
	return nil
}
