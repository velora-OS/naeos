package infrastructure

type Provider string

const (
	ProviderAWS   Provider = "aws"
	ProviderGCP   Provider = "gcp"
	ProviderAzure Provider = "azure"
	ProviderLocal Provider = "local"
)

type Infrastructure struct {
	Provider    Provider          `json:"provider,omitempty"`
	Region      string            `json:"region,omitempty"`
	Resources   []Resource        `json:"resources,omitempty"`
	Networking  []Network         `json:"networking,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

type Resource struct {
	Name string            `json:"name"`
	Kind string            `json:"kind,omitempty"`
	Spec map[string]string `json:"spec,omitempty"`
}

type Network struct {
	Name   string   `json:"name"`
	Kind   string   `json:"kind,omitempty"`
	Ports  []int    `json:"ports,omitempty"`
}
