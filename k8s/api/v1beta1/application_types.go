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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	acmeioutils "github.com/nathanbrophy/portfolio-demo/k8s/utils"
)

const (
	NAME            string = "acme-application"
	SERVICE_ACCOUNT string = NAME + "-sa"
	VERSION         string = "v1.0.0"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplicationBoilerPlate defines a set of boilder plate information that is useful across all reconciliation steps
type ApplicationBoilerPlate struct {
	// ServiceAccount is an optional flag to define the name of the service account to generate
	//+optional
	ServiceAccount *string `json:"serviceAccount,omitempty"`

	// ImagePullSecrets is an array of pull secrets to bind to the generated service account
	//+optional
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`

	// NamePrefix allows the resource name generation to be overriden, and can be derived when not present
	//+optional
	NamePrefix *string `json:"namePrefix,omitempty"`

	// Version defines the version for the static k8s labels
	//+optional
	Version *string `json:"version,omitempty"`
}

// ApplicationApplication defines information that is used to deploy the application itself and ensure it can run on the cluster environment
type ApplicationApplication struct {
	// Image defines the FQDN / Pull Location for the container image to run and is required
	Image *string `json:"image"`

	// Replicas is the number of replicas to run for the downstream deployment
	//+optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Port is the port to expose from the container
	//++optional
	Port *int32 `json:"port,omitempty"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Application defines the application specific information to use in reconciliation
	Application *ApplicationApplication `json:"application"`

	// BoilerPlate defines bootstrap / helpful information and metadata to be used and is not tied directly to the application
	BoilerPlate *ApplicationBoilerPlate `json:"boilerPlate,omitempty"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Progressing defines if the install is currently in progress or completed
	Progressing bool `json:"progressing"`

	// Reason defines why progressing is true or false
	Reason string `json:"reason"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

/**
 *
 * Star method API design function implementations
 * please reference the api.go package for complete documentation
 * on what each of the delta point functions completes.
 *
 */

func (a *Application) Replicas() *int32 {
	if a == nil || a.Spec.Application == nil || a.Spec.Application.Replicas == nil {
		return acmeioutils.Int32PointerGenerator(1)
	}

	return a.Spec.Application.Replicas
}

func (a *Application) Image() string {
	if a == nil || a.Spec.Application == nil {
		return ""
	}

	return *a.Spec.Application.Image
}

func (a *Application) Port() *int32 {
	if a == nil || a.Spec.Application == nil || a.Spec.Application.Port == nil {
		return acmeioutils.Int32PointerGenerator(8081)
	}

	return a.Spec.Application.Port
}

func (a *Application) ServiceAccount() *string {
	if a == nil || a.Spec.BoilerPlate == nil || a.Spec.BoilerPlate.ServiceAccount == nil {
		return acmeioutils.StringPointerGenerator(SERVICE_ACCOUNT)
	}

	return a.Spec.BoilerPlate.ServiceAccount
}

func (a *Application) ImagePullSecrets() []string {
	if a == nil || a.Spec.BoilerPlate == nil {
		return []string{}
	}

	return a.Spec.BoilerPlate.ImagePullSecrets
}

func (a *Application) Name() *string {
	if a == nil || a.Spec.BoilerPlate == nil || a.Spec.BoilerPlate.NamePrefix == nil {
		return acmeioutils.StringPointerGenerator(NAME)
	}

	return a.Spec.BoilerPlate.NamePrefix
}

func (a *Application) Version() *string {
	if a == nil || a.Spec.BoilerPlate == nil || a.Spec.BoilerPlate.Version == nil {
		return acmeioutils.StringPointerGenerator(VERSION)
	}

	return a.Spec.BoilerPlate.Version
}

func (a *Application) Instancer() *string {
	uuid := string(a.ObjectMeta.UID)
	truncMax := 6

	return acmeioutils.StringPointerGenerator(uuid[:truncMax])
}

// Reduired in order to interact with the control plane
func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
