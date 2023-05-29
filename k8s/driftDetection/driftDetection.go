package driftdetection

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DriftDetectionFunc defines the function family used to determine if drift exists on the cluster or not
// takes in two objects as input:
//
//	in: The object we have reconciled from the desired cluster state
//	out: The object that already exists on the cluster we must reconcile againse
//
// and perform a deep comparision on thetwo objects to determine if any such drift exists
type DriftDetectionFunc func(in, out client.Object) bool

// Deployment implements DriftDetectionFunc for the deployment resource
func Deployment(in, out client.Object) bool {
	lhs := in.(*appsv1.Deployment)
	rhs := out.(*appsv1.Deployment)

	drift := *lhs.Spec.Replicas != *rhs.Spec.Replicas
	drift = drift || !reflect.DeepEqual(lhs.Spec.Selector, rhs.Spec.Selector)
	drift = drift || !reflect.DeepEqual(lhs.Spec.Template.Labels, rhs.Spec.Template.Labels)
	drift = drift || !reflect.DeepEqual(lhs.Spec.Template.Spec.Containers, rhs.Spec.Template.Spec.Containers)

	return drift
}

// Service implements DriftDetectionFunc for the Service resource
func Service(in, out client.Object) bool {
	lhs := in.(*corev1.Service)
	rhs := out.(*corev1.Service)

	drift := !reflect.DeepEqual(lhs.Spec.Selector, rhs.Spec.Selector)
	drift = drift || !reflect.DeepEqual(lhs.Spec.Ports, rhs.Spec.Ports)

	return drift
}

// ServiceAccount implements DriftDetectionFunc for the ServiceAccount resource
func ServiceAccount(in, out client.Object) bool {
	lhs := in.(*corev1.ServiceAccount)
	rhs := out.(*corev1.ServiceAccount)

	return !reflect.DeepEqual(lhs.ImagePullSecrets, rhs.ImagePullSecrets)
}

// Ingress implements DriftDetectionFunc for the ServiceAccount resource
func Ingress(in, out client.Object) bool {
	lhs := in.(*networkingv1.Ingress)
	rhs := out.(*networkingv1.Ingress)

	return !reflect.DeepEqual(lhs.Spec.Rules, rhs.Spec.Rules)
}
