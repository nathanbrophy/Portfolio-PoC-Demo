package generators

import (
	"fmt"

	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	DefaultDeploymentGenerator     Generator = &DeploymentGeneratorV1{}
	DefaultServiceGenerator        Generator = &ServiceGeneratorV1{}
	DefaultServiceAccountGenerator Generator = &ServiceAccountGeneratorV1{}
)

type Generator interface {
	Object(acmeapi.Application) client.Object
}

func labelsGenerator(in acmeapi.Application) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       *in.Name(),
		"app.kubernetes.io/instance":   fmt.Sprintf("%s-%s", *in.Name(), *in.Instancer()),
		"app.kubernetes.io/version":    *in.Version(),
		"app.kubernetes.io/managed-by": "acme-controller",
		"app.kubernetes.io/part-of":    "acme-application",
	}
}

func generateAppSelector(in acmeapi.Application) *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": *in.Name(),
		},
	}
}

func generateContainerPorts(in acmeapi.Application) []corev1.ContainerPort {
	return []corev1.ContainerPort{
		{
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: *in.Port(),
		},
	}
}
