package api

type Application interface {
	Replicas() *int32
	ServiceAccount() *string
	ImagePullSecrets() []string
	Image() string
	Port() *int32
	Name() *string
	Version() *string
	Instancer() *string
}
