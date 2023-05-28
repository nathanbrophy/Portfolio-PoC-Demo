package generators

import (
	"reflect"
	"testing"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmetest "github.com/nathanbrophy/portfolio-demo/k8s/utils/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestServiceAccountGeneratorV1_Object(t *testing.T) {
	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		s    *ServiceAccountGeneratorV1
		args args
		want client.Object
	}{
		{
			name: "defaults",
			s:    &ServiceAccountGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServiceAccount",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "acme-application-sa",
					Labels: acmetest.DefaultMatchLabels(),
				},
				ImagePullSecrets: []corev1.LocalObjectReference{},
			},
		},
		{
			name: "no defaults",
			s:    &ServiceAccountGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServiceAccount",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "service-account-test-1",
					Labels: acmetest.NonDefaultMatchLabels(),
				},
				ImagePullSecrets: []corev1.LocalObjectReference{
					{
						Name: "docker.io",
					},
					{
						Name: "quay.io",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServiceAccountGeneratorV1{}
			if got := s.Object(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceAccountGeneratorV1.Object() = %v, want %v", got, tt.want)
			}
		})
	}
}
