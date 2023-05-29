package generators

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmeutils "github.com/nathanbrophy/portfolio-demo/k8s/utils"
)

const ingressClass string = "alb"

// IngressGeneratorV1 implemented the Generator interface for the service k8s manifest type
type IngressGeneratorV1 struct{}

// Object will generate the reconciled ingress from the expected cluster state
func (s *IngressGeneratorV1) Object(in acmeapi.Application) client.Object {
	pType := networkingv1.PathType("Prefix")
	generated := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   *in.Name(),
			Labels: labelsGenerator(in),
			Annotations: map[string]string{
				"alb.ingress.kubernetes.io/scheme":      "internet-facing",
				"alb.ingress.kubernetes.io/target-type": "ip",
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: acmeutils.StringPointerGenerator(ingressClass),
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: *in.Name(),
											Port: networkingv1.ServiceBackendPort{
												Number: *in.Port(),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return generated
}
