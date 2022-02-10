package lubeconfig

import (
	"fmt"
	"github.com/ertia-io/config/pkg/entities"
	"os"
	"time"
)

type ErtiaConfig struct {
	ErtiaConfig entities.LubeConfig
	ErtiaConfigPath string
}






func SpawnLubeConfig() (*entities.LubeConfig,error){
	lubeConfig := &entities.LubeConfig{
		Created: time.Now(),
		Updated: time.Now(),
		Nodes:[]entities.LubeNode{},
		Provider:"GLESYS", //TODO: Pluggable
		//Provider:"HETZNER", //TODO: Pluggable
	}
	err := lubeConfig.Persist()

	return lubeConfig, err
}

type LubeConfigOption func(config *entities.LubeConfig) *entities.LubeConfig

func WithProjectId(projectId string) LubeConfigOption {
	return func(cfg *entities.LubeConfig) *entities.LubeConfig{
		cfg.ProjectID = projectId
		return cfg
	}
}


func WithProviderToken(token string) LubeConfigOption {
	return func(cfg *entities.LubeConfig) *entities.LubeConfig{
		cfg.APIToken = token
		return cfg
	}
}

func Init(opts ...LubeConfigOption) (error){

	fmt.Printf("Creating directory: %s \n", config.LubePath())
	err := os.MkdirAll(config.LubePath(), 0777)
	if (err != nil) {
		return err
	}

	fmt.Printf("Creating directory: %s \n", config.LubeKeysPath())
	err = os.MkdirAll(config.LubeKeysPath(), 0777)
	if (err != nil) {
		return err
	}

	fmt.Printf("Creating directory: %s \n", config.LubeKubePath())
	err = os.MkdirAll(config.LubeKubePath(), 0777)
	if (err != nil) {
		return err
	}

	cfg, err := SpawnLubeConfig()
	if (err != nil) {
		return err
	}

	for _, opt := range opts {
		cfg = opt(cfg)
	}


	pKeys, err := keys.GetPublicKeys()
	if (err != nil) {
		return err
	}

	cfg.SSHKeys = pKeys

	err = cfg.Persist()
	if (err != nil) {
		return err
	}

	return nil
}