package generators

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	"github.com/nathanbrophy/portfolio-demo/k8s/utils"
)

// DeploymentGeneratorV1 implemented the Generator interface for the deployment k8s manifest type
type DeploymentGeneratorV1 struct{}

// Object will generate the reconciled service from the expected cluster state
func (d *DeploymentGeneratorV1) Object(in acmeapi.Application) client.Object {
	// The labels from generation need to be merged
	// so that the resulting Pod template has a valid
	// spec definition.
	baseLabels := labelsGenerator(in)
	selectorLabels := generateAppSelector(in)

	for k, v := range selectorLabels.MatchLabels {
		baseLabels[k] = v
	}

	selectorLabels.MatchLabels = baseLabels

	generated := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   *in.Name(),
			Labels: labelsGenerator(in),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: in.Replicas(),
			Selector: selectorLabels,
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
				},
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: baseLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "application-container",
							Image:           in.Image(),
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
							Ports:                    generateContainerPorts(in),
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
						},
					},
					ServiceAccountName:            *in.ServiceAccount(),
					TerminationGracePeriodSeconds: utils.Int64PointerGenerator(90),
				},
			},
		},
	}

	return generated
}
