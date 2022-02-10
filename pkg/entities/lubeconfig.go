package entities

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"lube/pkg/config"
	"lube/pkg/keys"
	"net"
	"os"
	"time"
)

type LubeConfig struct {
	ProjectID string `json:"project"`
	APIToken string `json:"apiToken"`
	SSHKeys []keys.LubeSSHKey `json:"sshKeys"`
	Nodes []LubeNode `json:"nodes"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Context string `json:"context"`
	Provider string `json:"provider"`
}

func(lc *LubeConfig) ToJSON() ([]byte, error) {
	return json.MarshalIndent(lc,"","  ")
}

func(lc *LubeConfig) WriteJSON(w io.Writer) (error) {
	enc := json.NewEncoder(w)
	enc.SetIndent("","  ")
	if err := enc.Encode(lc); err != nil {
		return err
	}

	return nil
}

func(lc *LubeConfig) Persist() error {
	file, err := os.OpenFile(config.LubeConfigPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if(err!=nil){
		return err
	}

	err = lc.WriteJSON(file)
	if(err!=nil){
		return err
	}
	return nil
}

func(lc *LubeConfig) FindNodeByIPV4(ip net.IP) (*LubeNode){

	for mi := range lc.Nodes {
		if(lc.Nodes[mi].IPV4.Equal(ip)){
			return &lc.Nodes[mi]
		}
	}
	return nil
}


func(lc *LubeConfig) FindNodeByID(id string) (*LubeNode){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ID == id){
			return &lc.Nodes[mi]
		}
	}
	return nil
}


func(lc *LubeConfig) FindNodeByName(name string) (*LubeNode){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].Name == name){
			return &lc.Nodes[mi]
		}
	}
	return nil
}

func(lc *LubeConfig) FindClusterMasterNode(clusterName string) (*LubeNode){
	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ClusterName == clusterName && lc.Nodes[mi].IsMaster){
			return &lc.Nodes[mi]
		}
	}
	return nil
}

func(lc *LubeConfig) FindKeyByID(id string) (*keys.LubeSSHKey){
	for mi := range lc.SSHKeys {
		if(lc.SSHKeys[mi].ID == id){
			return &lc.SSHKeys[mi]
		}
	}
	return nil
}


func (lc *LubeConfig) UpdateNode(node *LubeNode) (*LubeConfig, error) {

	node.Updated=time.Now()

	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ID == node.ID){
			lc.Nodes[mi] = *node
		}
	}
	err := lc.Persist()

	return lc, err
}


func (lc *LubeConfig) RemoveNode(node *LubeNode) (*LubeConfig, error) {

	for mi := range lc.Nodes {
		if(lc.Nodes[mi].ID == node.ID){
			lc.Nodes = removeNode(lc.Nodes, mi )
		}
	}
	err := lc.Persist()

	return lc, err
}



func (lc *LubeConfig) UpdateKey(key *keys.LubeSSHKey) (*LubeConfig, error) {

	for mi := range lc.SSHKeys {
		if(lc.SSHKeys[mi].ID == key.ID){
			lc.SSHKeys[mi] = *key
		}
	}
	err := lc.Persist()

	return lc, err
}

func (lc *LubeConfig) RemoveKey(key *keys.LubeSSHKey) (*LubeConfig, error) {

	for mi := range lc.SSHKeys {
		if(lc.SSHKeys[mi].ID == key.ID){
			lc.SSHKeys = removeKey(lc.SSHKeys, mi )
		}
	}
	err := lc.Persist()

	return lc, err
}


func ParseLubeConfig(path string) (*LubeConfig, error){
	contents, err := ioutil.ReadFile(path)
	if(err!=nil){
		return nil, err
	}

	lubeConfig := &LubeConfig{}

	err = json.Unmarshal(contents, lubeConfig)
	if(err!=nil){
		return nil, err
	}

	return lubeConfig, nil
}

func removeKey(s []keys.LubeSSHKey, i int) []keys.LubeSSHKey {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}


func removeNode(s []LubeNode, i int) []LubeNode {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}