package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/steven-zou/go-plugin/pkg/context"
	"github.com/steven-zou/go-plugin/pkg/plugin"
)

//Hello ...
type Hello struct {
	Name string
}

func main() {
	wk, err := os.Getwd()
	if err != nil {
		PrintError(err)
	}
	pluginBaseDir := filepath.Join(wk, "plugins")
	log.Printf("[INFO]: Plugin base dir: %s\n", pluginBaseDir)

	pluginManager := plugin.DefaultManager
	pluginManager.SetPluginBaseDir(pluginBaseDir)
	pluginManager.LoadPlugins()

	spec, executor, err := pluginManager.GetPlugin("sample")
	if err != nil {
		PrintError(err)
	}

	log.Printf("[INFO]: Plugin '%s' loaded with version '%s'\n", spec.Name, spec.Version)
	pContext := context.Background()
	pContext.SetValue("sample", &Hello{"hello go-plugin"})
	executor(pContext)
}

//PrintError ...
func PrintError(err error) {
	log.Printf("[Error]: %s\n", err)
	os.Exit(1)
}
