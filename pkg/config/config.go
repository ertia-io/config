package config

import (
	"github.com/containers/podman/pkg/util"
	"os"
	"path/filepath"
)

func ErtiaPath() string {
	var configPath string
	if path, ok := os.LookupEnv("ERTIADIR"); ok {
		configPath = path
	} else{
		homeDir, err := util.HomeDir()
		if(err!=nil){
			homeDir = "/opt/ertia/"
		}
		configPath = filepath.Join(homeDir, ".ertia", ErtiaContext())
	}
	return configPath
}

func ErtiaConfigPath() string {
	var configPath string
	if path, ok := os.LookupEnv("ERTIACONFIG"); ok {
		configPath = path
	} else{
		configPath = filepath.Join(ErtiaPath(),"/config.json")
	}
	return configPath
}

func ErtiaKeysPath() string {
	var configPath string
	if path, ok := os.LookupEnv("ERTIAKEYS"); ok {
		configPath = path
	} else{
		configPath = filepath.Join(ErtiaPath(),"/.ssh")
	}
	return configPath
}


func ErtiaKubePath() string {
	var configPath string
	if path, ok := os.LookupEnv("ERTIAKUBE"); ok {
		configPath = path
	} else{
		configPath = filepath.Join(ErtiaPath(),"/.kube")
	}
	return configPath
}

func ErtiaKubeConfigPath() string {
	return filepath.Join(ErtiaKubePath(),"/config")
}


func ErtiaContext() string {
	var ertiaContext string
	if lc, ok := os.LookupEnv("ERTIACONTEXT"); ok {
		ertiaContext = lc
	} else{
		ertiaContext = "DEFAULT"
	}
	return ertiaContext

}