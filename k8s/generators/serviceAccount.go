package generators

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
)

type ServiceAccountGeneratorV1 struct{}

func (s *ServiceAccountGeneratorV1) Object(in acmeapi.Application) client.Object {
	pullSecrets := in.ImagePullSecrets()
	lors := make([]corev1.LocalObjectReference, len(pullSecrets))
	for i, ps := range pullSecrets {
		lors[i] = corev1.LocalObjectReference{
			Name: ps,
		}
	}

	generated := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   *in.ServiceAccount(),
			Labels: labelsGenerator(in),
		},
		ImagePullSecrets: lors,
	}

	return generated
}
