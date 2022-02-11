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

type Node struct {
	ID          string           `json:"id"`
	ProviderID  string           `json:"providerId"`
	NodeToken  string           `json:"nodeToken"`
	IsMaster    bool             `json:"isMaster"`
	MasterIP    net.IP			 `json:"masterIP"`
	Tags            []string     `json:"tags"`
	Name            string       `json:"name"`
	IPV4            net.IP       `json:"ipv4"`
	IPV6            net.IP       `json:"ipv6"`
	Status          string       `json:"status"`
	Retries         int          `json:"retries"`
	Error           string       `json:"error"`
	Created         time.Time    `json:"created"`
	Updated         time.Time    `json:"updated"`
	Deleted         *time.Time   `json:"deleted"`
	Features        NodeFeatures `json:"features"`
	Dependencies    []Dependency `json:"dependencies"`
	Deployments     []Deployment `json:"deployments"`
	InstallPassword string       `json:"installPassword"`
	InstallUser     string       `json:"installUser"`
	//TODO: Add data as needed
}


type NodeFeatures map[string] bool

func(ma *Node) NeedsAdapting() bool{
	return ma.Status == NodeStatusNew || ma.Status == NodeStatusActive || ma.Status == NodeStatusRetrying
}
func(n *Node) Retry() (*Node){
	if(n.Retries> 10){
		n.Status = NodeStatusFailing
		return n
	}

	n.Retries = n.Retries+1
	n.Status = NodeStatusRetrying
	return n
}

func (ma *Node) Requires(dependency string) bool {
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

func (ma *Node) Fulfils(dependency string) bool {
	for i := range ma.Dependencies {
		if(ma.Dependencies[i].Name == dependency){
			if(ma.Dependencies[i].Status == DependencyStatusReady){
				return true
			}
		}
	}
	return false
}