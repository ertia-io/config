package entities

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"time"
)

const DefaultProvider = "GLESYS"
const ErrorNoReservableProjects = "Could not find any available projects to reserve"

type Project struct {
	ProjectID string          `json:"project"`
	Name string 				`json:"name"`
	APIToken string           `json:"apiToken"`
	SSHKey *SSHKey `json:"sshKey"`
	Nodes []Node              `json:"nodes"`
	Created time.Time         `json:"created"`
	Updated time.Time         `json:"updated"`
	Context string            `json:"context"`
	Provider string           `json:"provider"`
	Reserved bool 			`json:"reserved"`
	Domain string `json:"domain"`
	Tags []string `json:"tags"`
}

type Projects []Project


func(cfg *Project) WithProvider(provider string) *Project{
	cfg.Provider = provider
	return cfg
}

func(cfg *Project) WithName(name string) *Project{
	cfg.Name = name
	return cfg
}

func(cfg *Project) WithAPIToken(apiToken string) *Project{
	cfg.APIToken = apiToken
	return cfg
}


func(lc *Project) ToJSON() ([]byte, error) {
	return json.MarshalIndent(lc,"","  ")
}

func(lc *Project) WriteJSON(w io.Writer) (error) {
	enc := json.NewEncoder(w)
	enc.SetIndent("","  ")
	if err := enc.Encode(lc); err != nil {
		return err
	}

	return nil
}

func (p *Project) Reserve() *Project{
	p.Reserved = true
	return p
}

func (p *Project) HasTag(tag string) bool {
	for _,t := range p.Tags {
		if(t == tag ) {
			return true
		}
	}

	return false
}

func (p *Project) Tag(tags ...string) *Project {

NewTags:
	for _, newTag := range tags {
		for _, tag := range p.Tags {
			if(tag == newTag) {
				continue NewTags
			}
		}
		p.Tags = append(p.Tags, newTag)
	}

	return p

}

func(lc *Project) FindNodeByIPV4(ip net.IP) (*Node){

	for mi := range lc.Nodes {
		if(lc.Nodes[mi].IPV4.Equal(ip)){
			return &lc.Nodes[mi]
		}
	}
	return nil
}


func(lc *Project) FindNodeByID(id string) (*Node){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ID == id){
			return &lc.Nodes[mi]
		}
	}
	return nil
}


func(lc *Project) FindNodeByName(name string) (*Node){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].Name == name){
			return &lc.Nodes[mi]
		}
	}
	return nil
}

func(lc *Project) FindClusterMasterNode(clusterName string) (*Node){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ClusterName == clusterName && lc.Nodes[mi].IsMaster){
			return &lc.Nodes[mi]
		}
	}
	return nil
}



func (lc *Project) UpdateNode(node *Node) (*Project) {

	node.Updated=time.Now()

	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ID == node.ID){
			lc.Nodes[mi] = *node
		}
	}

	return lc
}


func (lc *Project) RemoveNode(node *Node) (*Project) {

	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ID == node.ID){
			lc.Nodes = removeNode(lc.Nodes, mi )
		}
	}
	return lc
}



func (lc *Project) UpdateKey(key *SSHKey) (*Project) {

	lc.SSHKey = key

	return lc
}

func (lc *Project) RemoveKey(key *SSHKey) (*Project) {

	lc.SSHKey = nil
	return lc
}


func ParseConfig(jsonStr string) (*Project, error){
	config := &Project{}

	err := json.Unmarshal([]byte(jsonStr), config)
	if(err!=nil){
		return nil, err
	}

	return config, nil
}

func ParseProjects(jsonStr string) (Projects, error){
	config := []Project{}

	err := json.Unmarshal([]byte(jsonStr), &config)
	if(err!=nil){
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

func WithID(projectId string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.ProjectID = projectId
		return cfg
	}
}


func WithName(name string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.Name = name
		return cfg
	}
}

func WithDomain(domain string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.Domain = domain
		return cfg
	}
}



func WithProviderToken(token string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.APIToken = token
		return cfg
	}
}

func NewProject(opts ...ProjectOption) (*Project, error){

	project := &Project{
		Created: time.Now(),
		Updated: time.Now(),
		Nodes:[]Node{},
		Provider:DefaultProvider, //TODO: Pluggable
	}

	for _, opt := range opts {
		project = opt(project)
	}


	pKey,_, err := GetPublicKeys()
	if (err != nil) {
		return nil,err
	}

	project.SSHKey = pKey

	if (err != nil) {
		return nil,err
	}

	return project,nil
}

func(ps Projects) ReserveOne() (Projects, *Project, error){
	for i, p := range ps {
		if(p.Reserved == false && !p.HasTag("ertia-pool-master")){
			ps[i].Reserved = true
			return ps, &p, nil
		}
	}
	return ps, nil, errors.New(ErrorNoReservableProjects)
}

