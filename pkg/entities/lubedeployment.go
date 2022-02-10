package entities

const (
	DeploymentStatusNew = "NEW"
	DeploymentStatusDeploying = "DEPLOYING"
	DeploymentStatusFailing = "FAILING"
	DeploymentStatusRetrying = "READY"
	DeploymentStatusReady = "READY"
)

type LubeDeployment struct {
	Name string `json:"name"`
	Status string `json:"status"`
	Retries int32 `json:"retries"`
}