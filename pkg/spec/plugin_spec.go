package spec

//Plugin is the corresponding structure of the 'plugin.json',
//which describe the basic metadata of the plugin
type Plugin struct {
	//Name of the plugin, required
	Name string

	//SemVer 2 version, required
	Version string

	//A one line sentence about the function of the plugin, optional
	Description string

	//The source repository of the plugin code, optional
	Home string

	//The maintainer list with 'Maintainer <email>' format of the plugin
	Maintainers []string

	//The source for loading the plugin with specified mode
	Source *Source

	//The HTTP service should be served by the plugin
	HTTPServices *HTTPServices
}

//Source defines the loading mode of the plugin
type Source struct {
	//The loading mode of the plugin
	//Support 'local_so', 'remote_git'
	Mode string

	//The path of the local so file or the URL of the remote git
	Path string
}

//HTTPServiceRoute defines the http/rest service endpoint served by the plugin
type HTTPServiceRoute struct {
	//The service endpoint
	Route string

	//The mothod of the service
	Method string

	//The label will be set into the plugin context to
	//let the plugin aware what kind of request is incoming for serving
	Label string
}

//HTTPServices defines the metadata of http service served by the plugin
type HTTPServices struct {
	//The http service provider
	//Support 'beego'
	Driver string

	//Routes should be enabled on the driver
	Routes []HTTPServiceRoute
}
