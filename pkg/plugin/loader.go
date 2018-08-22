package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	sys_plugin "plugin"

	"github.com/steven-zou/go-plugin/pkg"

	"github.com/steven-zou/go-plugin/pkg/context"
	"github.com/steven-zou/go-plugin/pkg/spec"
)

//Loader defines the plugin load flow
type Loader interface {
	//Scan the plugin base dir and get the plugin candidates
	Scan(pluginBaseDir string) ([]string, error)

	//Parse the plugin metadata
	Parse(pluginPath string) (*spec.Plugin, error)

	//Load the plugin executor
	Load(plugin *spec.Plugin) (spec.PluginExecutor, error)
}

//BaseLoader is an default implementation of Loader interface
type BaseLoader struct{}

//Scan implements same method of Loader interface
func (bl *BaseLoader) Scan(pluginBaseDir string) ([]string, error) {
	if len(pluginBaseDir) == 0 {
		return nil, errors.New("plugin base dir path is empty")
	}

	if !pkg.FileExists(pluginBaseDir) {
		return nil, fmt.Errorf("plugin base dir '%s' is not existing", pluginBaseDir)
	}

	files, err := ioutil.ReadDir(pluginBaseDir)
	if err != nil {
		return nil, err
	}

	candidates := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			candidates = append(candidates, filepath.Join(pluginBaseDir, f.Name()))
		}
	}

	return candidates, nil
}

//Parse implements same method of Loader interface
func (bl *BaseLoader) Parse(pluginPath string) (*spec.Plugin, error) {
	return nil, errors.New("not implemented in BaseLoader, leverage the validator to get the plugin spec")
}

//Load implements same method of Loader interface
func (bl *BaseLoader) Load(plugin *spec.Plugin) (spec.PluginExecutor, error) {
	if plugin == nil {
		return nil, errors.New("nil plugin spec")
	}

	if !pkg.FileExists(plugin.Source.Path) {
		return nil, fmt.Errorf("plugin so file '%s' is not existsing", plugin.Source.Path)
	}

	p, err := sys_plugin.Open(plugin.Source.Path)
	if err != nil {
		return nil, err
	}

	exec, err := p.Lookup("Execute")
	if err != nil {
		return nil, err
	}

	pExec, ok := exec.(func(ctx context.PluginContext) error)
	if !ok {
		return nil, fmt.Errorf("failed to lookup entry function 'Execute' in plugin so file '%s'", plugin.Source.Path)
	}

	return pExec, nil
}
