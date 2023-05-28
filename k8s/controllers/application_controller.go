/*
Copyright 2023 Nathan Brophy.

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
	"reflect"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	acmeiov1beta1 "github.com/nathanbrophy/portfolio-demo/k8s/api/v1beta1"
	acmegdrift "github.com/nathanbrophy/portfolio-demo/k8s/driftDetection"
	acmegenerators "github.com/nathanbrophy/portfolio-demo/k8s/generators"
)

type ReconcileWrapper struct {
	Driftor      acmegdrift.DriftDetectionFunc
	Manifest     client.Object
	ObjectLoader client.Object
}

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func gvk(obj client.Object) schema.GroupVersionKind {
	return obj.GetObjectKind().GroupVersionKind()
}

func (r *ApplicationReconciler) updateStatus(
	logger logr.Logger,
	ctx context.Context,
	req ctrl.Request,
	progressing bool,
	err error,
) error {
	found := &acmeiov1beta1.Application{}
	_ = r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, found)

	newStatus := found.Status.DeepCopy()

	newStatus.Progressing = progressing
	newStatus.Reason = "reconciling cluster state"

	if !progressing {
		newStatus.Reason = "completed"
	}

	if err != nil {
		newStatus.Progressing = false
		newStatus.Reason = fmt.Sprintf("failed to reconcile cluster state due to error: %v", err)
	}

	// A deep equal reflection is required to prevent an infinite reconciliation loop from occuring.
	//
	// Since the status is managed as a subresource of the API we are required to access it
	// through the subresource .Status() mutator to propogate changes.
	if !reflect.DeepEqual(found.Status, *newStatus) {
		found.Status = *newStatus
		if err := r.Status().Update(ctx, found); err != nil {
			if errors.IsTooManyRequests(err) || errors.IsConflict(err) {
				return nil
			}
			return err
		}
	}

	return nil
}

//+kubebuilder:rbac:groups=acme.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=acme.io,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=acme.io,resources=applications/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services;serviceaccounts,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reconcileLogger := log.FromContext(ctx)

	reconcileLogger.Info("begin reconciliation")
	reconcileLogger.Info("attempting to retrieve the CR instance from current namespace")

	// Load the CR instances to watch from the namespace, so that the cluster state can
	// be reconciled to the specification defined in the CR.
	cr := &acmeiov1beta1.Application{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, cr); err != nil {
		if errors.IsNotFound(err) {
			// The object was deleted or no longer exists, do not infinitely loop
			// attempting to load the resource after deletion
			return ctrl.Result{}, nil
		}
		reconcileLogger.Error(err, "cannot get CR from namespace, requeueing attempt and trying again")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	// Define a collection of information required to reconcile cluster state
	toReconcile := []ReconcileWrapper{
		{
			Driftor:      acmegdrift.Deployment,
			Manifest:     acmegenerators.DefaultDeploymentGenerator.Object(cr),
			ObjectLoader: &appsv1.Deployment{},
		},
		{
			Driftor:      acmegdrift.Service,
			Manifest:     acmegenerators.DefaultServiceGenerator.Object(cr),
			ObjectLoader: &corev1.Service{},
		},
		{
			Driftor:      acmegdrift.ServiceAccount,
			Manifest:     acmegenerators.DefaultServiceAccountGenerator.Object(cr),
			ObjectLoader: &corev1.ServiceAccount{},
		},
	}

	if err := r.updateStatus(reconcileLogger, ctx, req, true, nil); err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	for _, reconcilers := range toReconcile {
		// Set the namespace for the generated manifest to
		// the namespace for the reconciling CR
		reconcilers.Manifest.SetNamespace(cr.Namespace)

		// If the controller reference is not set, then things like
		// cascading object garbage collection, state enforcement, and
		// ownership are not propogated correctly.
		ctrl.SetControllerReference(cr, reconcilers.Manifest, r.Scheme)
		objGVK := gvk(reconcilers.Manifest)
		reconcileLogger.Info(
			"attempting to reconcile a manifest to correct cluster state for given CR in context",
			"group",
			objGVK.Group,
			"kind",
			objGVK.Kind,
			"version",
			objGVK.Version,
		)

		// Attempt to create the object, assume it does no exist, and
		// handle edge scenarios from there.
		if err := r.Client.Create(ctx, reconcilers.Manifest); err != nil {
			if errors.IsAlreadyExists(err) {
				// In this scenario the object already exists on the cluster,
				// so the reconciler must check to see if the object has
				// drifted from the defined cluster state, or if there is no
				// detectable drift.
				//
				// Without checking the cluster state drift, the reconciler
				// will be stuck in an infinite loop as it will perform no-op
				// udpdates to the existing cluster objects.
				found := reconcilers.ObjectLoader
				found.SetNamespace(reconcilers.Manifest.GetNamespace())
				found.SetName(reconcilers.Manifest.GetName())
				if err := r.Client.Get(ctx, client.ObjectKeyFromObject(found), found); err != nil {
					// We cannot determine if drift exists or not if we cannot
					// grab the current object state from the cluster.
					if err := r.updateStatus(reconcileLogger, ctx, req, false, err); err != nil {
						return ctrl.Result{RequeueAfter: time.Second * 5}, err
					}
					return ctrl.Result{RequeueAfter: time.Second * 5}, err
				}
				if !reconcilers.Driftor(reconcilers.Manifest, found) {
					// No drift detected is an indicator that
					// no reconciliation is required.
					return ctrl.Result{}, nil
				}

				reconcileLogger.Info("found a conflicting object state on the cluster, overriding definition to match expected cluster state")
				if err := r.Client.Update(ctx, reconcilers.Manifest); err != nil {
					// When this happens the cluster is in a dirty state where
					// there is drift that cannot be recovered from, meaning the
					// current cluster state is not valid to the CR definition
					reconcileLogger.Error(err, "unable to update object to restore expected cluster state")
					if err := r.updateStatus(reconcileLogger, ctx, req, false, err); err != nil {
						return ctrl.Result{RequeueAfter: time.Second * 5}, err
					}
					return ctrl.Result{RequeueAfter: time.Second * 5}, err
				}
			} else {
				reconcileLogger.Error(err, "unable to create require downstream manifest to support application deployment")
				if err := r.updateStatus(reconcileLogger, ctx, req, false, err); err != nil {
					return ctrl.Result{RequeueAfter: time.Second * 5}, err
				}
				return ctrl.Result{RequeueAfter: time.Second * 5}, err
			}
		}
	}

	if err := r.updateStatus(reconcileLogger, ctx, req, false, nil); err != nil {
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&acmeiov1beta1.Application{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
