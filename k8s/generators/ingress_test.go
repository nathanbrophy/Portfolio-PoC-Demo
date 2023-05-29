package generators

import (
	"reflect"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmeioutils "github.com/nathanbrophy/portfolio-demo/k8s/utils"
	acmetest "github.com/nathanbrophy/portfolio-demo/k8s/utils/test"
)

func TestIngressGeneratorV1_Object(t *testing.T) {
	pType := networkingv1.PathType("Prefix")

	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		s    *IngressGeneratorV1
		args args
		want client.Object
	}{
		{
			name: "defaults",
			s:    &IngressGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: &networkingv1.Ingress{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "acme-application",
					Labels: acmetest.DefaultMatchLabels(),
					Annotations: map[string]string{
						"alb.ingress.kubernetes.io/scheme":      "internet-facing",
						"alb.ingress.kubernetes.io/target-type": "ip",
					},
				},
				Spec: networkingv1.IngressSpec{
					IngressClassName: acmeioutils.StringPointerGenerator("alb"),
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
													Name: "acme-application",
													Port: networkingv1.ServiceBackendPort{
														Number: *acmeioutils.Int32PointerGenerator(8081),
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
			},
		},
		{
			name: "no defaults",
			s:    &IngressGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: &networkingv1.Ingress{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Ingress",
					APIVersion: "networking.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "example-prefix",
					Labels: acmetest.NonDefaultMatchLabels(),
					Annotations: map[string]string{
						"alb.ingress.kubernetes.io/scheme":      "internet-facing",
						"alb.ingress.kubernetes.io/target-type": "ip",
					},
				},
				Spec: networkingv1.IngressSpec{
					IngressClassName: acmeioutils.StringPointerGenerator("alb"),
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
													Name: "example-prefix",
													Port: networkingv1.ServiceBackendPort{
														Number: *acmeioutils.Int32PointerGenerator(5051),
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IngressGeneratorV1{}
			if got := s.Object(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IngressGeneratorV1.Object() = %v, want %v", got, tt.want)
			}
		})
	}
}
