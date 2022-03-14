package entities

const (
	ProvisioningStatusNew       = "NEW"
	ProvisioningStatusDeploying = "DEPLOYING"
	ProvisioningStatusFailing   = "FAILING"
	ProvisioningStatusRetrying  = "RETRYING"
	ProvisioningStatusReady     = "READY"
)

type Provisioning struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	Retries      int    `json:"retries"`
	Organization string `json:"organization"`
}
