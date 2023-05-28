package test_utils

import (
	acmeapi "github.com/nathanbrophy/portfolio-demo/k8s/api"
	acmeiov1beta1 "github.com/nathanbrophy/portfolio-demo/k8s/api/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func GenerateCRWithDefaults() acmeapi.Application {
	generated := &acmeiov1beta1.Application{
		ObjectMeta: v1.ObjectMeta{
			UID: types.UID("12345678"),
		},
	}
	generated.Spec.Application = &acmeiov1beta1.ApplicationApplication{
		Image: func(x string) *string { return &x }("example.com/test-image:v1.0"),
	}

	return generated
}

func GenerateCRWithNoDefaults() acmeapi.Application {
	generated := &acmeiov1beta1.Application{
		ObjectMeta: v1.ObjectMeta{
			UID: types.UID("12345678"),
		},
		Spec: acmeiov1beta1.ApplicationSpec{
			Application: &acmeiov1beta1.ApplicationApplication{
				Image:    func(x string) *string { return &x }("example.com/test-image:v1.0"),
				Replicas: func(x int32) *int32 { return &x }(5),
				Port:     func(x int32) *int32 { return &x }(5051),
			},
			BoilerPlate: &acmeiov1beta1.ApplicationBoilerPlate{
				ServiceAccount: func(x string) *string { return &x }("service-account-test-1"),
				ImagePullSecrets: []string{
					"docker.io",
					"quay.io",
				},
				NamePrefix: func(x string) *string { return &x }("example-prefix"),
				Version:    func(x string) *string { return &x }("v1.0.1"),
			},
		},
	}

	return generated
}

func DefaultMatchLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "acme-application",
		"app.kubernetes.io/instance":   "acme-application-123456",
		"app.kubernetes.io/version":    "v1.0.0",
		"app.kubernetes.io/managed-by": "acme-controller",
		"app.kubernetes.io/part-of":    "acme-application",
	}
}

func NonDefaultMatchLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "example-prefix",
		"app.kubernetes.io/instance":   "example-prefix-123456",
		"app.kubernetes.io/version":    "v1.0.1",
		"app.kubernetes.io/managed-by": "acme-controller",
		"app.kubernetes.io/part-of":    "acme-application",
	}
}
