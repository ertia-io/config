package entities

import (
	"net"
	"time"
)



const (
	NodeStatusNew        = "NEW"
	NodeStatusDeploying  = "DEPLOYING"
	NodeStatusActive     = "ACTIVE"
	NodeStatusReady      = "READY"
	NodeStatusFailing    = "FAILING"
	NodeStatusRetrying   = "RETRYING"
	NodeStatusRestarting = "RESTARTING"
	NodeStatusStopped    = "STOPPED"
	NodeStatusError      = "ERROR"
	NodeStatusDeleted    = "DELETED"
)

type LubeNode struct {
	ID          string           `json:"id"`
	ClusterName string           `json:"clusterName"`
	ProviderID  string           `json:"providerId"`
	NodeToken  string           `json:"nodeToken"`
	IsMaster    bool             `json:"isMaster"`
	MasterIP    net.IP			 `json:"masterIP"`
	Tags        []string         `json:"tags"`
	Name        string           `json:"name"`
	IPV4        net.IP           `json:"ipv4"`
	IPV6        net.IP           `json:"ipv6"`
	Status      string           `json:"status"`
	Retries     int              `json:"retries"`
	Error	    string				`json:"error"`
	Created     time.Time        `json:"created"`
	Updated     time.Time        `json:"updated"`
	Deleted     *time.Time       `json:"deleted"`
	Features    LubeNodeFeatures `json:"features"`
	Dependencies []LubeDependency `json:"dependencies"`
	Deployments []LubeDeployment `json:"deployments"`
	InstallPassword string `json:"installPassword"`
	InstallUser string `json:"installUser"`
	//TODO: Add data as needed
}


type LubeNodeFeatures map[string] bool

func(ma *LubeNode) NeedsAdapting() bool{
	return ma.Status == NodeStatusNew || ma.Status == NodeStatusActive || ma.Status == NodeStatusRetrying
}
func(ma *LubeNode) Retry() (*LubeNode, error){
	if(ma.Retries> 10){
		ma.Status = NodeStatusFailing
		ma, err := ma.Persist()

		if(err!=nil){
			return ma, err
		}
		return ma, nil
	}

	ma.Retries = ma.Retries+1
	ma.Status = NodeStatusRetrying
	ma, err := ma.Persist()

	if(err!=nil){
		return ma, err
	}

	return ma, nil

}

func (ma *LubeNode) Persist() (*LubeNode,error) {
	cfg, err := ParseLubeConfig(config.LubeConfigPath())
	if(err!=nil){
		return ma, err
	}


	ma.Updated = time.Now()
	_, err = cfg.UpdateNode(ma)

	if(err!=nil){
		return ma, err
	}

	return ma, nil
}

func (ma *LubeNode) Requires(dependency string) bool {
	for i := range ma.Dependencies {
		if(ma.Dependencies[i].Name == dependency){
			if(ma.Dependencies[i].Status == DependencyStatusNew ||
				ma.Dependencies[i].Status == DependencyStatusRetrying ||
				ma.Dependencies[i].Status == DependencyStatusWaiting ){
				return true
			}
		}
	}
	return false
}

func (ma *LubeNode) Fulfils(dependency string) bool {
	for i := range ma.Dependencies {
		if(ma.Dependencies[i].Name == dependency){
			if(ma.Dependencies[i].Status == DependencyStatusReady){
				return true
			}
		}
	}
	return false
}