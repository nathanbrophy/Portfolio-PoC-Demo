/*
Copyright 2023 Nathan Brophy.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	acmeioutils "github.com/nathanbrophy/portfolio-demo/k8s/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestApplication_Replicas(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *int32
	}{
		{
			name:   "default",
			fields: fields{},
			want:   acmeioutils.Int32PointerGenerator(1),
		},
		{
			name: "no default",
			fields: fields{
				Spec: ApplicationSpec{
					Application: &ApplicationApplication{
						Replicas: acmeioutils.Int32PointerGenerator(3),
					},
				},
			},
			want: acmeioutils.Int32PointerGenerator(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.Replicas(); *got != *tt.want {
				t.Errorf("Application.Replicas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Image(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "base scenario",
			fields: fields{
				Spec: ApplicationSpec{
					Application: &ApplicationApplication{
						Image: acmeioutils.StringPointerGenerator("example.io/image/v1.0"),
					},
				},
			},
			want: "example.io/image/v1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.Image(); got != tt.want {
				t.Errorf("Application.Image() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Port(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *int32
	}{
		{
			name:   "default",
			fields: fields{},
			want:   acmeioutils.Int32PointerGenerator(8081),
		},
		{
			name: "no default",
			fields: fields{
				Spec: ApplicationSpec{
					Application: &ApplicationApplication{
						Port: acmeioutils.Int32PointerGenerator(5051),
					},
				},
			},
			want: acmeioutils.Int32PointerGenerator(5051),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.Port(); *got != *tt.want {
				t.Errorf("Application.Port() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_ServiceAccount(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{
			name:   "default",
			fields: fields{},
			want:   acmeioutils.StringPointerGenerator("acme-application-sa"),
		},
		{
			name: "no default",
			fields: fields{
				Spec: ApplicationSpec{
					BoilerPlate: &ApplicationBoilerPlate{
						ServiceAccount: acmeioutils.StringPointerGenerator("test-sa"),
					},
				},
			},
			want: acmeioutils.StringPointerGenerator("test-sa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.ServiceAccount(); *got != *tt.want {
				t.Errorf("Application.ServiceAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_ImagePullSecrets(t *testing.T) {
	basefmt := "docker-%d-.io"
	var td [5]string

	for i := 0; i < 5; i++ {
		td[i] = fmt.Sprintf(basefmt, rand.Int31())
	}
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "default",
			fields: fields{},
			want:   []string{},
		},
		{
			name: "no default (1)",
			fields: fields{
				Spec: ApplicationSpec{
					BoilerPlate: &ApplicationBoilerPlate{
						ImagePullSecrets: []string{"docker.io"},
					},
				},
			},
			want: []string{"docker.io"},
		},
		{
			name: "no default (2)",
			fields: fields{
				Spec: ApplicationSpec{
					BoilerPlate: &ApplicationBoilerPlate{
						ImagePullSecrets: []string{"docker.io", "quay.io"},
					},
				},
			},
			want: []string{"docker.io", "quay.io"},
		},
		{
			name: "no default fuzz",
			fields: fields{
				Spec: ApplicationSpec{
					BoilerPlate: &ApplicationBoilerPlate{
						ImagePullSecrets: td[:],
					},
				},
			},
			want: td[:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.ImagePullSecrets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Application.ImagePullSecrets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Name(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.Name(); got != tt.want {
				t.Errorf("Application.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_Version(t *testing.T) {
	noDefault := fmt.Sprintf("v%d.%d.%d", rand.Int31n(10), rand.Int31n(100), rand.Int31n(1000))
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{
			name:   "default",
			fields: fields{},
			want:   acmeioutils.StringPointerGenerator("v1.0.0"),
		},
		{
			name: "no default",
			fields: fields{
				Spec: ApplicationSpec{
					BoilerPlate: &ApplicationBoilerPlate{
						Version: acmeioutils.StringPointerGenerator(noDefault),
					},
				},
			},
			want: acmeioutils.StringPointerGenerator(noDefault),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.Version(); *got != *tt.want {
				t.Errorf("Application.Version() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestApplication_Instancer(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       ApplicationSpec
		Status     ApplicationStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   *string
	}{
		{
			name: "default",
			fields: fields{
				ObjectMeta: metav1.ObjectMeta{
					UID: types.UID("12345678-asdg"),
				},
			},
			want: acmeioutils.StringPointerGenerator("123456"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Application{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := a.Instancer(); *got != *tt.want {
				t.Errorf("Application.Instancer() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
