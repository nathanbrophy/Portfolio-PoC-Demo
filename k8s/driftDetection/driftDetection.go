package driftdetection

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DriftDetectionFunc func(in, out client.Object) bool

func Deployment(in, out client.Object) bool {
	lhs := in.(*appsv1.Deployment)
	rhs := out.(*appsv1.Deployment)

	drift := *lhs.Spec.Replicas != *rhs.Spec.Replicas
	drift = drift || !reflect.DeepEqual(lhs.Spec.Selector, rhs.Spec.Selector)
	drift = drift || !reflect.DeepEqual(lhs.Spec.Template.Labels, rhs.Spec.Template.Labels)
	drift = drift || !reflect.DeepEqual(lhs.Spec.Template.Spec.Containers, rhs.Spec.Template.Spec.Containers)

	return drift
}

func Service(in, out client.Object) bool {
	lhs := in.(*corev1.Service)
	rhs := out.(*corev1.Service)

	drift := !reflect.DeepEqual(lhs.Spec.Selector, rhs.Spec.Selector)
	drift = drift || !reflect.DeepEqual(lhs.Spec.Ports, rhs.Spec.Ports)

	return drift
}

func ServiceAccount(in, out client.Object) bool {
	lhs := in.(*corev1.ServiceAccount)
	rhs := out.(*corev1.ServiceAccount)

	return !reflect.DeepEqual(lhs.ImagePullSecrets, rhs.ImagePullSecrets)
}
