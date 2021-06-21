/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	resourcev1alpha1 "github.com/dweber019/k8s-copier/api/v1alpha1"

	"github.com/jinzhu/copier"
)

// CopyResourceReconciler reconciles a CopyResource object
type CopyResourceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=resource.w3tec.ch,resources=copyresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resource.w3tec.ch,resources=copyresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=resource.w3tec.ch,resources=copyresources/finalizers,verbs=update
// +kubebuilder:rbac:groups=v1,resources=secret,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CopyResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *CopyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("copyresource", req.NamespacedName)

	copyResource := &resourcev1alpha1.CopyResource{}
	err := r.Get(ctx, req.NamespacedName, copyResource)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("CopyResource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get CopyResource.")
		return ctrl.Result{}, err
	}

	namespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      copyResource.Spec.MetaName,
	}
	sourceResource, _ := StringToStruct(copyResource.Spec.Kind)
	err = r.Client.Get(ctx, namespacedName, sourceResource)
	if err != nil && !errors.IsNotFound(err) {
		log.Error(err, "Secret error.")
		return ctrl.Result{Requeue: true}, nil
	}

	targetResource, _ := StringToStruct(copyResource.Spec.Kind)
	targetResource, _ = cloneResource(copyResource.Spec.Kind, sourceResource, targetResource)
	targetResource.SetResourceVersion("")
	targetResource.SetUID("")
	targetResource.SetNamespace(copyResource.Spec.TargetNamespace)
	targetResource.SetName(copyResource.Namespace + "-" + copyResource.Name)

	targetNamespacedName := types.NamespacedName{
		Namespace: targetResource.GetNamespace(),
		Name:      targetResource.GetName(),
	}
	targetNamespacedObject, _ := StringToStruct(copyResource.Spec.Kind)
	err = r.Client.Get(ctx, targetNamespacedName, targetNamespacedObject)

	if copyResource.Status.ResourceVersion == "" || copyResource.Status.ResourceVersion != sourceResource.GetResourceVersion() || errors.IsNotFound(err) {
		err = r.Client.Create(ctx, targetResource)
		if err != nil && errors.IsAlreadyExists(err) {
			err = r.Client.Update(ctx, targetResource)
		}
		if err == nil {
			copyResource.Status.ResourceVersion = sourceResource.GetResourceVersion()
			err := r.Status().Update(ctx, copyResource)
			if err != nil {
				log.Error(err, "Failed to update CopyResource status")
				return ctrl.Result{Requeue: true}, err
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CopyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&resourcev1alpha1.CopyResource{}).
		Complete(r)
}

func StringToStruct(kind string) (client.Object, error) {
	switch kind {
	case "Secret":
		return &v1.Secret{}, nil
	case "ConfigMap":
		return &v1.ConfigMap{}, nil
	default:
		return nil, fmt.Errorf("%s is not a known resource name", kind)
	}
}

func cloneResource(kind string, source client.Object, target client.Object) (client.Object, error) {
	switch kind {
	case "Secret":
		copier.Copy(target.(*v1.Secret), source.(*v1.Secret))
		return target, nil
	case "ConfigMap":
		copier.Copy(target.(*v1.ConfigMap), source.(*v1.ConfigMap))
		return target, nil
	default:
		return nil, fmt.Errorf("%s is not a known resource name", kind)
	}
}
