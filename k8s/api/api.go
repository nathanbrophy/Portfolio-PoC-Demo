package api

// Application defines the interface that all versions of the API must adhere to in the star API versioning scheme
type Application interface {
	// Replicas returns the number of replicas a deployment must have
	Replicas() *int32

	// ServiceAccount is an optional field which will pass in a predefined service account name to use
	ServiceAccount() *string

	// ImagePullSecrets defines a set a image pul secrets to bind to the service account
	ImagePullSecrets() []string

	// Image defines the FQDN for the pull location for the Application's container image
	Image() string

	// Port defines the port to expose from the Application's container
	Port() *int32

	// Name defines the name of the overall Application suite and can be derrived from existing information
	Name() *string

	// Version defines the application version running
	Version() *string

	// Instancer derives the UUID instance truncation from the CR's generated UUID in etcd
	Instancer() *string
}
