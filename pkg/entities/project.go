package entities

import (
	"encoding/json"
	"errors"
	"github.com/lucasepe/codename"
	"github.com/segmentio/ksuid"
	"io"
	"net"
	"time"
)

const DefaultProvider = "GLESYS"
const ErrorNoReservableProjects = "Could not find any available projects to reserve"

type Project struct {
	ID string `json:"id"`
	Provider string           `json:"provider"`
	ProviderID string `json:"providerID"` // Provider Project ID
	ProviderToken string `json:"providerToken"`
	Name string 				`json:"name"`
	SSHKey *SSHKey `json:"sshKey"`
	Nodes []Node              `json:"nodes"`
	Created time.Time         `json:"created"`
	Updated time.Time         `json:"updated"`
	Context string            `json:"context"`

	Reserved bool 			`json:"reserved"`
	Domain string `json:"domain"`
	Tags []string `json:"tags"`
}

type Projects []Project


func(cfg *Project) WithProvider(provider string) *Project{
	cfg.Provider = provider
	return cfg
}

func(cfg *Project) WithProviderID(providerProjectId string) *Project{
	cfg.ProviderID = providerProjectId
	return cfg
}

func(cfg *Project) WithProviderToken(token string) *Project{
	cfg.ProviderToken = token
	return cfg
}

func(cfg *Project) WithName(name string) *Project{
	cfg.Name = name
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

func(lc *Project) FindMasterNode() (*Node){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].IsMaster){
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

func WithID(id string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.ID = id
		return cfg
	}
}



func WithProvider(provider string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.Provider = provider
		return cfg
	}
}


func WithProviderID(projectProviderId string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.ProviderID = projectProviderId
		return cfg
	}
}


func WithProviderToken(token string) ProjectOption {
	return func(cfg *Project) *Project{
		cfg.ProviderToken = token
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


func NodeName() string{
	rng, err := codename.DefaultRNG()
	if err != nil {
		return ""
	}
	return codename.Generate(rng, 4)

}

func (cfg *Project) AddNode() (*Project, *Node, error){

	thisIsMaster := true

	master := cfg.FindMasterNode()

	var masterIp net.IP
	var nodeToken string
	if(master != nil){
		thisIsMaster = false
		masterIp  = master.IPV4
		nodeToken = master.NodeToken
	}

	newNode := Node{
		ID:       ksuid.New().String(),
		Name:     NodeName(),
		IsMaster: thisIsMaster,
		MasterIP: masterIp,
		NodeToken: nodeToken,
		Status:   "NEW",
		Created:  time.Now(),
		Updated:  time.Now(),
		Deleted:  nil,
		Features: NodeFeatures{},
	}

	cfg.Nodes = append(cfg.Nodes, newNode)

	return cfg, &newNode, nil
}