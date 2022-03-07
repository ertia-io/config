package entities

const (
	DeploymentStatusNew       = "NEW"
	DeploymentStatusDeploying = "DEPLOYING"
	DeploymentStatusFailing   = "FAILING"
	DeploymentStatusRetrying  = "RETRYING"
	DeploymentStatusReady     = "READY"
)

type Deployment struct {
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	Tag     string `json:"tag"`
	Token   string `json:"token"`
	Status  string `json:"status"`
	Retries int32  `json:"retries"`
}
