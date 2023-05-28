package generators

import (
	"encoding/json"
	"reflect"
	"testing"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmetest "github.com/nathanbrophy/portfolio-demo/k8s/utils/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestServiceGeneratorV1_Object(t *testing.T) {
	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		s    *ServiceGeneratorV1
		args args
		want client.Object
	}{
		{
			name: "defaults",
			s:    &ServiceGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "acme-application",
					Labels: acmetest.DefaultMatchLabels(),
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"app": "acme-application",
					},
					Ports: []corev1.ServicePort{
						{
							Protocol:   corev1.ProtocolTCP,
							Port:       8081,
							TargetPort: intstr.FromInt(8081),
						},
					},
				},
			},
		},
		{
			name: "no defaults",
			s:    &ServiceGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "example-prefix",
					Labels: acmetest.NonDefaultMatchLabels(),
				},
				Spec: corev1.ServiceSpec{
					Selector: map[string]string{
						"app": "example-prefix",
					},
					Ports: []corev1.ServicePort{
						{
							Protocol:   corev1.ProtocolTCP,
							Port:       5051,
							TargetPort: intstr.FromInt(5051),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceGeneratorV1{}
			if got := s.Object(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(tt.want, "", "  ")
				t.Errorf("ServiceGeneratorV1.Object() = \n%v\n, want =\n%v\n", string(g), string(w))
			}
		})
	}
}
