package pkg

import (
	"os"
)

const (
	//PluginJSONFileName is the pre-defined filename of plugin metadata json file
	PluginJSONFileName = "plugin.json"

	//PluginSourceModeLocal defines the local mode
	PluginSourceModeLocal = "local_so"

	//PluginSourceModeRemote defines the remote mode
	PluginSourceModeRemote = "remote_git"
)

//FileExists check the existence of the specified file
//If file exists, return true
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

//IsDir checks if the file is a dir
func IsDir(filePath string) bool {
	fi, err := os.Stat(filePath)

	return err == nil && fi.Mode().IsDir()
}
