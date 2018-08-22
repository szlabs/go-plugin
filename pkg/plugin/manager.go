package plugin

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/steven-zou/go-plugin/pkg"
	"github.com/steven-zou/go-plugin/pkg/spec"
)

//DefaultManager is the default plugin manager
var DefaultManager = NewBaseManager()

//Manager defines the related operations of one plugin manager
//should support.
//Manager is used to load, organize and maintain the plugins.
type Manager interface {
	//Set the base dir where to load plugins.
	//If the dir does not exist or it's not a dir
	//an error will be returned.
	SetPluginBaseDir(dir string) error

	//Load all the plugins from the base plugin dir.
	//Any issues happened, an error will be returned.
	LoadPlugins() error

	//Load plugin with the specified name.
	//If failed to load, an error will be returned.
	LoadPlugin(name string) error

	//Unload plugin with the specified name.
	//If failed to unload, an error will be returned.
	UnloadPlugin(name string) error

	//Get the plugin with the specified name.
	//If plugin is not existing, an error will be returned.
	GetPlugin(name string) (*spec.Plugin, spec.PluginExecutor, error)
}

//BaseManager is implemented as default plugin manager
type BaseManager struct {
	//Keep the base dir of the plugins
	basePluginBaseDir string

	//The plugin loader
	loader Loader

	//The plugin validator
	validtor Validator

	//The list to keep the loaded one
	store Store
}

//NewBaseManager is constructor of BaseManager
func NewBaseManager() Manager {
	return &BaseManager{
		loader: &BaseLoader{},
		validtor: NewBaseValidatorChain(
			&JSONFileValidator{},
			&SpecValidator{},
			&LocalSourceValidator{}),
		store: NewBaseStore(),
	}
}

//SetPluginBaseDir implements the interface method
func (bm *BaseManager) SetPluginBaseDir(dir string) error {
	if len(dir) > 0 {
		if pkg.FileExists(dir) {
			bm.basePluginBaseDir = dir
			return nil
		}
	}

	return fmt.Errorf("%s is not a valid plugin base dir path", dir)
}

//LoadPlugins implements the interface method
func (bm *BaseManager) LoadPlugins() error {
	//scan plugin base dir
	paths, err := bm.loader.Scan(bm.basePluginBaseDir)
	if err != nil {
		return err
	}

	if len(paths) == 0 {
		//No plugin dirs
		return nil
	}

	//loop all
	for _, p := range paths {
		log.Printf("[INFO]: Found plugin: %s\n", p)
		if err := bm.loadPlugin(p); err != nil {
			log.Printf("[ERROR]: Plugin loading error: %s\n", err)
		}
	}

	log.Printf("[INFO]: %d plugins loaded", bm.store.Size())

	return nil
}

//LoadPlugin implements the interface method
func (bm *BaseManager) LoadPlugin(name string) error {
	if len(name) == 0 {
		return errors.New("plugin name cannot be empty")
	}

	pluginPath := filepath.Join(bm.basePluginBaseDir, name)

	return bm.loadPlugin(pluginPath)
}

//UnloadPlugin implements the interface method
func (bm *BaseManager) UnloadPlugin(name string) error {
	if len(name) == 0 {
		return errors.New("plugin name cannot be empty")
	}

	if _, ok := bm.store.Get(name); !ok {
		return fmt.Errorf("plugin %s is not loaded", name)
	}

	if _, ok := bm.store.Remove(name); !ok {
		return fmt.Errorf("failed to unload plugin %s", name)
	}

	return nil
}

//GetPlugin implements the interface method
func (bm *BaseManager) GetPlugin(name string) (*spec.Plugin, spec.PluginExecutor, error) {
	if len(name) == 0 {
		return nil, nil, errors.New("plugin name cannot be empty")
	}

	pluginItem, ok := bm.store.Get(name)
	if !ok {
		return nil, nil, fmt.Errorf("plugin with name '%s' is not existing", name)
	}

	return pluginItem.Spec, pluginItem.Executor, nil
}

func (bm *BaseManager) loadPlugin(pluginPath string) error {
	//validate
	validateRes, err := bm.validtor.Validate(pluginPath)
	if err != nil {
		log.Printf("[INFO]: Valdate plugin [FAILED]: %s", pluginPath)
		return err
	}
	log.Printf("[INFO]: Valdate plugin [SUCCESS]: %s", pluginPath)

	//Convert validate result to plugin spec object
	pluginSpec, ok := validateRes.(*spec.Plugin)
	if !ok {
		log.Println("[ERROR]: Failed to convert validation result to plugin spec")
		return errors.New("Failed to convert validation result to plugin spec")
	}
	//load
	exec, err := bm.loader.Load(pluginSpec)
	if err != nil {
		log.Printf("[INFO]: Load plugin [FAILED]: %s:%s", pluginSpec.Name, pluginSpec.Version)
		return err
	}
	log.Printf("[INFO]: Load plugin [SUCCESS]: %s:%s", pluginSpec.Name, pluginSpec.Version)

	//Save
	bm.store.Put(&spec.PluginItem{
		Spec:     pluginSpec,
		Executor: exec,
	}, true)

	return nil
}
