package entities

const (
	DeploymentStatusNew = "NEW"
	DeploymentStatusDeploying = "DEPLOYING"
	DeploymentStatusFailing = "FAILING"
	DeploymentStatusRetrying = "READY"
	DeploymentStatusReady = "READY"
)

type Deployment struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Status string `json:"status"`
	Retries int32 `json:"retries"`
}