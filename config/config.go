package config

import (
	"os"
)

func LubePath() string {
	var lubeConfigPath string
	if path, ok := os.LookupEnv("LUBEDIR"); ok {
		lubeConfigPath = path
	} else{
		lubeConfigPath = "/home/robin/.lube" +"/" + LubeContext()
	}
	return lubeConfigPath
}

func LubeConfigPath() string {
	var lubeConfigPath string
	if path, ok := os.LookupEnv("LUBECONFIG"); ok {
		lubeConfigPath = path
	} else{
		lubeConfigPath = LubePath() + "/config.json"
	}
	return lubeConfigPath
}

func LubeKeysPath() string {
	var lubeConfigPath string
	if path, ok := os.LookupEnv("LUBEKEYS"); ok {
		lubeConfigPath = path
	} else{
		lubeConfigPath = LubePath() +"/.ssh"
	}
	return lubeConfigPath
}


func LubeKubePath() string {
	var lubeConfigPath string
	if path, ok := os.LookupEnv("LUBEKUBE"); ok {
		lubeConfigPath = path
	} else{
		lubeConfigPath = LubePath() + "/.kube"
	}
	return lubeConfigPath
}

func LubeKubeConfigPath() string {
	return LubeKubePath()+"/config"
}


func LubeContext() string {
	var lubeContext string
	if lc, ok := os.LookupEnv("LUBECONTEXT"); ok {
		lubeContext = lc
	} else{
		lubeContext = "DEFAULT"
	}
	return lubeContext

}