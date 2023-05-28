package generators

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
)

// ServiceGeneratorV1 implemented the Generator interface for the service k8s manifest type
type ServiceGeneratorV1 struct{}

// Object will generate the reconciled service from the expected cluster state
func (s *ServiceGeneratorV1) Object(in acmeapi.Application) client.Object {
	generated := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   *in.Name(),
			Labels: labelsGenerator(in),
		},
		Spec: corev1.ServiceSpec{
			Selector: generateAppSelector(in).MatchLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       *in.Port(),
					TargetPort: intstr.FromInt(int(*in.Port())),
				},
			},
		},
	}

	return generated
}
