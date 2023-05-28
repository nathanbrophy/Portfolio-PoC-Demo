package generators

import (
	"reflect"
	"testing"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmetest "github.com/nathanbrophy/portfolio-demo/k8s/utils/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_generateContainerPorts(t *testing.T) {
	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		args args
		want []corev1.ContainerPort
	}{
		{
			name: "defaults",
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: []corev1.ContainerPort{
				{
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: 8081,
				},
			},
		},
		{
			name: "no defaults",
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: []corev1.ContainerPort{
				{
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: 5051,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateContainerPorts(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateContainerPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateAppSelector(t *testing.T) {
	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		args args
		want *metav1.LabelSelector
	}{
		{
			name: "defaults",
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "acme-application",
				},
			},
		},
		{
			name: "no defaults",
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "example-prefix",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateAppSelector(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateAppSelector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_labelsGenerator(t *testing.T) {
	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "defaults",
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: acmetest.DefaultMatchLabels(),
		},
		{
			name: "no defaults",
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: acmetest.NonDefaultMatchLabels(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := labelsGenerator(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("labelsGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
