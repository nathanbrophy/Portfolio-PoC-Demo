package driftdetection

import (
	"math/rand"
	"testing"

	acmeioutils "github.com/nathanbrophy/portfolio-demo/k8s/utils"
	acmetest "github.com/nathanbrophy/portfolio-demo/k8s/utils/test"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestDeployment(t *testing.T) {
	defaultLabelsSel := acmetest.DefaultMatchLabels()
	defaultLabelsSel["app"] = "acme-application"
	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "acme-application",
			Labels: acmetest.DefaultMatchLabels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(x int32) *int32 { return &x }(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: defaultLabelsSel,
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
				},
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: defaultLabelsSel,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "application-container",
							Image:           "example.com/test-image:v1.0",
							ImagePullPolicy: corev1.PullAlways,
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.LifecycleHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"sleep 30",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8081,
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
					},
					ServiceAccountName:            "acme-application-sa",
					TerminationGracePeriodSeconds: func(x int64) *int64 { return &x }(90),
				},
			},
		},
	}

	dDiff := d.DeepCopy()
	dDiff.Spec.Replicas = acmeioutils.Int32PointerGenerator(rand.Int31n(15) + 2)

	type args struct {
		in  client.Object
		out client.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "defaults match",
			args: args{
				in:  d,
				out: d,
			},
			want: false,
		},
		{
			name: "defaults do not match",
			args: args{
				in:  d,
				out: dDiff,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Deployment(tt.args.in, tt.args.out); got != tt.want {
				t.Errorf("Deployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService(t *testing.T) {
	s := &corev1.Service{
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
	}

	sDiff := s.DeepCopy()
	sDiff.Spec.Selector["app"] = "changed"

	type args struct {
		in  client.Object
		out client.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "defaults match",
			args: args{
				in:  s,
				out: s,
			},
			want: false,
		},
		{
			name: "defaults do not match",
			args: args{
				in:  s,
				out: sDiff,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Service(tt.args.in, tt.args.out); got != tt.want {
				t.Errorf("Service() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceAccount(t *testing.T) {
	s := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "acme-application-sa",
			Labels: acmetest.DefaultMatchLabels(),
		},
		ImagePullSecrets: []corev1.LocalObjectReference{},
	}

	sDiff := s.DeepCopy()
	sDiff.ImagePullSecrets = []corev1.LocalObjectReference{
		{
			Name: "local-secret",
		},
	}

	type args struct {
		in  client.Object
		out client.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "defaults match",
			args: args{
				in:  s,
				out: s,
			},
			want: false,
		},
		{
			name: "defaults do not match",
			args: args{
				in:  s,
				out: sDiff,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ServiceAccount(tt.args.in, tt.args.out); got != tt.want {
				t.Errorf("ServiceAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIngress(t *testing.T) {
	pType := networkingv1.PathType("Prefix")
	generated := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "example",
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
											Name: "example",
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
	}

	copy := generated.DeepCopy()
	copy.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Number = *acmeioutils.Int32PointerGenerator(9091)

	type args struct {
		in  client.Object
		out client.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "deafult match",
			args: args{
				in:  generated,
				out: generated,
			},
			want: false,
		},
		{
			name: "deafult do not match",
			args: args{
				in:  generated,
				out: copy,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ingress(tt.args.in, tt.args.out); got != tt.want {
				t.Errorf("Ingress() = %v, want %v", got, tt.want)
			}
		})
	}
}
