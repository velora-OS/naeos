package cloud

import (
	"fmt"
	"time"
)

type CloudProvider string

const (
	AWS   CloudProvider = "aws"
	GCP   CloudProvider = "gcp"
	Azure CloudProvider = "azure"
)

// ResourceTypes maps abstract resource types to supported kinds.
const (
	ResourceStorage     = "storage"
	ResourceCompute     = "compute"
	ResourceDatabase    = "database"
	ResourceCache       = "cache"
	ResourceQueue       = "queue"
	ResourceCDN         = "cdn"
	ResourceServerless  = "serverless"
	ResourceMonitoring  = "monitoring"
	ResourceSecrets     = "secrets"
	ResourceDNS         = "dns"
	ResourceNetworking  = "networking"
)

var SupportedResourceTypes = []string{
	ResourceStorage,
	ResourceCompute,
	ResourceDatabase,
	ResourceCache,
	ResourceQueue,
	ResourceCDN,
	ResourceServerless,
	ResourceMonitoring,
	ResourceSecrets,
	ResourceDNS,
	ResourceNetworking,
}

type DeployConfig struct {
	Provider    CloudProvider
	Region      string
	Project     string
	Environment string
	Resources   []Resource
}

type Resource struct {
	Name string
	Type string
	Spec map[string]interface{}
}

type DeployResult struct {
	Provider   CloudProvider
	Resources  []DeployedResource
	Terraform  string
	Status     string
	Timestamp  time.Time
}

type DeployedResource struct {
	Name string
	Type string
	ID   string
	ARN  string
}

type PlanResult struct {
	Resources     []Resource   `json:"resources"`
	CostEstimate  CostEstimate `json:"cost_estimate"`
}

type CloudAdapter interface {
	Name() string
	Provider() CloudProvider
	Validate(config *DeployConfig) error
	Plan(config *DeployConfig) (*PlanResult, error)
	Deploy(config *DeployConfig) (*DeployResult, error)
	Destroy(config *DeployConfig) error
	ExportTerraform(config *DeployConfig) (string, error)
}

var adapterCache = map[CloudProvider]CloudAdapter{
	AWS:   &AWSAdapter{},
	GCP:   &GCPAdapter{},
	Azure: &AzureAdapter{},
}

func GetAdapter(provider CloudProvider) (CloudAdapter, error) {
	if adapter, ok := adapterCache[provider]; ok {
		return adapter, nil
	}
	return nil, fmt.Errorf("unsupported provider: %s", provider)
}
