package spec

//PluginItem is composite of plugin spec and executor
type PluginItem struct {
	//Plugin spec
	Spec *Plugin

	//Plugin executor
	Executor PluginExecutor
}
