package generators

import (
	"encoding/json"
	"reflect"
	"testing"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmetest "github.com/nathanbrophy/portfolio-demo/k8s/utils/test"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestDeploymentGeneratorV1_Object(t *testing.T) {
	defaultLabelsSel := acmetest.DefaultMatchLabels()
	defaultLabelsSel["app"] = "acme-application"

	noDefaultLabelsSel := acmetest.NonDefaultMatchLabels()
	noDefaultLabelsSel["app"] = "example-prefix"

	type args struct {
		in acmeapi.Application
	}
	tests := []struct {
		name string
		d    *DeploymentGeneratorV1
		args args
		want client.Object
	}{
		{
			name: "defaults",
			d:    &DeploymentGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithDefaults(),
			},
			want: &appsv1.Deployment{
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
			},
		},
		{
			name: "no defaults",
			d:    &DeploymentGeneratorV1{},
			args: args{
				in: acmetest.GenerateCRWithNoDefaults(),
			},
			want: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:   "example-prefix",
					Labels: acmetest.NonDefaultMatchLabels(),
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: func(x int32) *int32 { return &x }(5),
					Selector: &metav1.LabelSelector{
						MatchLabels: noDefaultLabelsSel,
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
							Labels: noDefaultLabelsSel,
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
											ContainerPort: 5051,
										},
									},
									TerminationMessagePath:   "/dev/termination-log",
									TerminationMessagePolicy: corev1.TerminationMessageReadFile,
								},
							},
							ServiceAccountName:            "service-account-test-1",
							TerminationGracePeriodSeconds: func(x int64) *int64 { return &x }(90),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DeploymentGeneratorV1{}
			if got := d.Object(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(tt.want, "", "  ")
				t.Errorf("DeploymentGeneratorV1.Object() = \n%v\n, want =\n%v\n", string(g), string(w))
			}
		})
	}
}
