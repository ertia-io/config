package entities

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"time"

	"github.com/lucasepe/codename"
	"github.com/segmentio/ksuid"
)

const DefaultProvider = "GLESYS"
const ErrorNoReservableProjects = "Could not find any available projects to reserve"

const ReserveGraceTime = time.Hour * 4

type Project struct {
	ID            string       `json:"id"`
	Provider      string       `json:"provider"`
	ProviderID    string       `json:"providerID"` // Provider Project ID
	ProviderToken string       `json:"providerToken"`
	Name          string       `json:"name"`
	SSHKey        *SSHKey      `json:"sshKey"`
	Nodes         []Node       `json:"nodes"`
	Created       time.Time    `json:"created"`
	Updated       time.Time    `json:"updated"`
	Context       string       `json:"context"`
	Deployments   []Deployment `json:"deployments"`
	Reserved      bool         `json:"reserved"`
	Delete        *time.Time   `json:"delete"`
	Domain        string       `json:"domain"`
	Tags          []string     `json:"tags"`
}

type Projects []Project

func (p *Project) WithProvider(provider string) *Project {
	p.Provider = provider
	return p
}

func (p *Project) WithProviderID(providerProjectId string) *Project {
	p.ProviderID = providerProjectId
	return p
}

func (p *Project) WithProviderToken(token string) *Project {
	p.ProviderToken = token
	return p
}

func (p *Project) WithName(name string) *Project {
	p.Name = name
	return p
}

func (p *Project) WithDeployments(deployments ...Deployment) *Project {
	p.Deployments = deployments
	return p
}

func (p *Project) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

func (p *Project) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if err := enc.Encode(p); err != nil {
		return err
	}

	return nil
}

func (p *Project) Reserve() *Project {
	p.Reserved = true
	deleteTime := time.Now().Add(ReserveGraceTime)
	p.Delete = &deleteTime
	return p
}

func (p *Project) HasTag(tag string) bool {
	for _, t := range p.Tags {
		if t == tag {
			return true
		}
	}

	return false
}

func (p *Project) Tag(tags ...string) *Project {
NewTags:
	for _, newTag := range tags {
		for _, tag := range p.Tags {
			if tag == newTag {
				continue NewTags
			}
		}
		p.Tags = append(p.Tags, newTag)
	}

	return p
}

func (p *Project) FindNodeByIPV4(ip net.IP) *Node {
	for mi := range p.Nodes {
		if p.Nodes[mi].IPV4.Equal(ip) {
			return &p.Nodes[mi]
		}
	}

	return nil
}

func (p *Project) FindNodeByID(id string) *Node {
	for mi := range p.Nodes {
		if p.Nodes[mi].ID == id {
			return &p.Nodes[mi]
		}
	}

	return nil
}

func (p *Project) FindNodeByName(name string) *Node {
	for mi := range p.Nodes {
		if p.Nodes[mi].Name == name {
			return &p.Nodes[mi]
		}
	}

	return nil
}

func (p *Project) FindMasterNode() *Node {
	for mi := range p.Nodes {
		if p.Nodes[mi].IsMaster {
			return &p.Nodes[mi]
		}
	}

	return nil
}

func (p *Project) UpdateNode(node *Node) *Project {
	node.Updated = time.Now()

	for mi := range p.Nodes {
		if p.Nodes[mi].ID == node.ID {
			p.Nodes[mi] = *node
		}
	}

	return p
}

func (p *Project) RemoveNode(node *Node) *Project {
	for mi := range p.Nodes {
		if p.Nodes[mi].ID == node.ID {
			p.Nodes = removeNode(p.Nodes, mi)
		}
	}

	return p
}

func (p *Project) UpdateKey(key *SSHKey) *Project {
	p.SSHKey = key
	return p
}

func (p *Project) RemoveKey(key *SSHKey) *Project {
	p.SSHKey = nil
	return p
}

func ParseConfig(jsonStr string) (*Project, error) {
	config := &Project{}

	err := json.Unmarshal([]byte(jsonStr), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseProjects(jsonStr string) (Projects, error) {
	config := []Project{}

	err := json.Unmarshal([]byte(jsonStr), &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func removeKey(s []SSHKey, i int) []SSHKey {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeNode(s []Node, i int) []Node {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type ProjectOption func(config *Project) *Project

func WithID(id string) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.ID = id
		return cfg
	}
}

func WithProvider(provider string) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.Provider = provider
		return cfg
	}
}

func WithProviderID(projectProviderId string) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.ProviderID = projectProviderId
		return cfg
	}
}

func WithProviderToken(token string) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.ProviderToken = token
		return cfg
	}
}

func WithName(name string) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.Name = name
		return cfg
	}
}

func WithDomain(domain string) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.Domain = domain
		return cfg
	}
}

func WithDeployments(deployments ...Deployment) ProjectOption {
	return func(cfg *Project) *Project {
		cfg.Deployments = deployments
		return cfg
	}
}

func NewProject(opts ...ProjectOption) (*Project, error) {
	project := &Project{
		Created:  time.Now(),
		Updated:  time.Now(),
		Nodes:    []Node{},
		Provider: DefaultProvider, //TODO: Pluggable
	}

	for _, opt := range opts {
		project = opt(project)
	}

	key, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	project.SSHKey = key

	return project, nil
}

func (ps Projects) ReserveOne() (Projects, *Project, error) {
	for i, p := range ps {
		if p.Reserved == false && p.Delete == nil {
			ps[i].Reserved = true
			deleteTime := time.Now().Add(ReserveGraceTime)
			ps[i].Delete = &deleteTime
			return ps, &p, nil
		}
	}
	return ps, nil, errors.New(ErrorNoReservableProjects)
}

func NodeName() string {
	rng, err := codename.DefaultRNG()
	if err != nil {
		return ""
	}
	return codename.Generate(rng, 4)

}

func (p *Project) AddNode() (*Project, *Node, error) {
	thisIsMaster := true
	master := p.FindMasterNode()

	var masterIp net.IP
	var nodeToken string

	if master != nil {
		thisIsMaster = false
		masterIp = master.IPV4
		nodeToken = master.NodeToken
	}

	newNode := Node{
		ID:        ksuid.New().String(),
		Name:      NodeName(),
		IsMaster:  thisIsMaster,
		MasterIP:  masterIp,
		NodeToken: nodeToken,
		Status:    NodeStatusNew,
		Created:   time.Now(),
		Updated:   time.Now(),
		Deleted:   nil,
		Features:  NodeFeatures{},
	}

	p.Nodes = append(p.Nodes, newNode)

	return p, &newNode, nil
}
