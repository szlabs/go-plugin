package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/steven-zou/go-plugin/pkg"
	"github.com/steven-zou/go-plugin/pkg/spec"
)

//Validator defines the behaviors of a plugin validator
type Validator interface {
	//Do validation with the provided params.
	//If meet any issues, an error will be returned.
	//If succeed, output the result which depends on the implementations.
	Validate(params ...interface{}) (interface{}, error)
}

//JSONFileValidator validates the existence of plugin.json.
type JSONFileValidator struct{}

//Validate is the implementation of Validator interface
func (jfv *JSONFileValidator) Validate(params ...interface{}) (interface{}, error) {
	if len(params) == 0 {
		return nil, errors.New("The plugin dir path is required")
	}

	pluginDirPath := fmt.Sprintf("%s", params[0])

	if !pkg.FileExists(pluginDirPath) {
		return nil, fmt.Errorf("File %s is not existsing", pluginDirPath)
	}

	if !pkg.IsDir(pluginDirPath) {
		return nil, fmt.Errorf("File %s is not a dir", pluginDirPath)
	}

	pluginJSONFile := filepath.Join(pluginDirPath, pkg.PluginJSONFileName)
	if !pkg.FileExists(pluginJSONFile) {
		return nil, fmt.Errorf("%s is not found under plugin dir %s", pkg.PluginJSONFileName, pluginDirPath)
	}

	data, err := ioutil.ReadFile(pluginJSONFile)
	if err != nil {
		return nil, err
	}

	//Load plugin.json
	pluginSpec := &spec.Plugin{}
	if err := json.Unmarshal(data, pluginSpec); err != nil {
		return nil, err
	}

	//Plugin dir name should be equal with the name of the plugin
	fi, err := os.Stat(pluginDirPath)
	if err != nil {
		//Actually, should not come here
		return nil, err
	}
	if fi.Name() != pluginSpec.Name {
		return nil, fmt.Errorf("Name conflicts: expect %s but got %s in the metadata json file", fi.Name(), pluginSpec.Name)
	}

	return pluginSpec, nil
}

//SpecValidator validates the plugin spec.
type SpecValidator struct{}

//Validate is the implementation of Validator interface
func (sv *SpecValidator) Validate(params ...interface{}) (interface{}, error) {
	if len(params) == 0 {
		return nil, errors.New("plugin json object is missing")
	}

	pluginSpec, ok := params[0].(*spec.Plugin)
	if !ok {
		return nil, errors.New("invalid plugin spec object")
	}

	if len(pluginSpec.Name) == 0 {
		return nil, errors.New("missing plugin name")
	}

	_, err := semver.NewVersion(pluginSpec.Version)
	if err != nil {
		return nil, err
	}

	if pluginSpec.Source == nil {
		return nil, errors.New("plugin source missing")
	}

	if pluginSpec.Source.Mode != pkg.PluginSourceModeLocal &&
		pluginSpec.Source.Mode != pkg.PluginSourceModeRemote {
		return nil, fmt.Errorf("Only support mode [%s, %s]", pkg.PluginSourceModeLocal, pkg.PluginSourceModeRemote)
	}

	return pluginSpec, nil
}

//LocalSourceValidator validate the local source
type LocalSourceValidator struct{}

//Validate is the implementation of Validator interface
func (lsv *LocalSourceValidator) Validate(params ...interface{}) (interface{}, error) {
	if len(params) < 2 {
		return nil, errors.New("plugin json object and plugin base dir are required")
	}

	pluginSpec, ok := params[0].(*spec.Plugin)
	if !ok {
		return nil, errors.New("invalid plugin spec object")
	}

	if pluginSpec.Source == nil {
		return nil, errors.New("plugin source missing")
	}

	//If the mode is not local mode, just ignore it
	if pluginSpec.Source.Mode != pkg.PluginSourceModeLocal {
		return pluginSpec, nil
	}

	pluginBaseDir := fmt.Sprintf("%s", params[1])
	//plugin so file path
	var pluginSoFilePath string
	if filepath.IsAbs(pluginSpec.Source.Path) {
		pluginSoFilePath = pluginSpec.Source.Path
	} else {
		pluginSoFilePath = filepath.Join(pluginBaseDir, pluginSpec.Source.Path)
	}

	if !pkg.FileExists(pluginSoFilePath) {
		return nil, fmt.Errorf("plugin so file %s is not existing", pluginSoFilePath)
	}

	if filepath.Ext(pluginSoFilePath) != ".so" {
		return nil, fmt.Errorf("%s.so file is missing", pluginSpec.Name)
	}

	//Override the so file path to absolute path
	pluginSpec.Source.Path = pluginSoFilePath

	return pluginSpec, nil
}

//RemoteSourceValidator validates the remote source
type RemoteSourceValidator struct{}

//Validate is the implementation of Validator interface
//TODO:
func (rsv *RemoteSourceValidator) Validate(params ...interface{}) (interface{}, error) {
	return nil, nil
}

//BaseValidatorChain build a validation pipeline with 'JSONFileValidator' and 'SpecValidator'.
type BaseValidatorChain struct {
	//The validator list
	validators []Validator
}

//NewBaseValidatorChain creates a validator chain
func NewBaseValidatorChain(validators ...Validator) Validator {
	bvc := &BaseValidatorChain{
		validators: make([]Validator, 0),
	}

	if len(validators) > 0 {
		bvc.validators = append(bvc.validators, validators...)
	}

	return bvc
}

//Validate is the implementation of Validator interface
func (bvc *BaseValidatorChain) Validate(params ...interface{}) (interface{}, error) {
	if len(bvc.validators) == 0 {
		return nil, errors.New("no validators")
	}
	if len(params) == 0 {
		return nil, errors.New("missing params")
	}

	var (
		result interface{}
		err    error
	)
	for _, vl := range bvc.validators {
		if result == nil {
			//The first validator
			result, err = vl.Validate(params...)
		} else {
			extendedParams := []interface{}{result}
			extendedParams = append(extendedParams, params...)
			result, err = vl.Validate(extendedParams...)
		}

		if err != nil {
			return nil, err
		}
	}

	//Please be aware that, the result can be nothing
	return result, nil
}
