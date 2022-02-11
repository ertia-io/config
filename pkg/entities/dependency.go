package entities

const (
	DependencyStatusNew = "NEW"
	DependencyStatusDeploying = "DEPLOYING"
	DependencyStatusFailing = "FAILING"
	DependencyStatusRetrying = "RETRYING"
	DependencyStatusReady = "READY"
	DependencyStatusWaiting = "WAITING"
)

type Dependency struct {
	Name string `json:"name"`
	Status string `json:"status"`
	Retries int32 `json:"retries"`
}