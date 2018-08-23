# go-plugin

go-plugin is intent on providing a plugin framework based on the [plugin](https://golang.org/pkg/plugin/) feature originated from golang 1.8. go-plugin works in the plugin consumer side to handle related stages of running a plugin, including discover, build, load, upgrade, manage, call and unload.

**If you want to do contributions to this project, WELCOME!**

## Maintainers
* Steven Zou <loneghost1982@gmail.com> (Originator)
* Maybe YOU

## Plugin specification

The plugin should follow the specifications below then it can be recognized by go-plugin framework.

### Entry method

The plugin should have a single unified and exported entry method with name `Execute`. The full signature is listed below:

```go
func Execute(ctx context.PluginContext) error {
    return nil
}
```

The entry method has a plugin context argument which extends the golang context interface and provides extra value operation methods. The detailed declarations of this context interface are shown below:

```go
//ValueContext defines behaviors of handling values with string keys
type ValueContext interface {
    //Get value by the key
    GetValue(key string) interface{}

    //Set value to context
    SetValue(key string, value interface{})
}

//PluginContext help to provide related information/parameters to the
//plugin execution entry method.
//PluginContext inherits all from the context.Context
type PluginContext interface {
    context.Context
    ValueContext
}

```

go-plugin provides a base plugin context for using.

```go
import "github.com/steven-zou/go-plugin/pkg/context"

context.BasePluginContext
```

### Metadata json file

go-plugin use a json file `plugin.json` to define and describe the plugin metadata. An example:

```json
{
    "name": "go-plugin",
    "version": "0.1.0",
    "description": "A go plugin management tool",
    "maintainers": ["szou@vmware.com"],
    "home": "https://github.com/steven-zou/go-plugin.git",
    "source": {
        "mode": "local_so",
        "path": "sample.so"
    },
    "http_services": {
        "driver": "beego",
        "routes": [
            { "route": "/api/plugin/sample/:id", "method": "GET", "label": "plugin.get" },
            { "route": "/api/plugin/samples", "method": "POST", "label": "plugin.post" }
        ]
    }
}
```

The spec details:

|        Field         |      Description       |      Required     |   Supported |
|----------------------|------------------------|-------------------|-------------|
|        name          | Name of the plugin     |        Y          |   Y         |
|        version       | Plugin version with semver 2 |  Y          |   Y         |
|        description   | One sentence to describe the plugin |  N   |   Y         |
|        maintainers   | A list of mails of maintainers |  N        |   Y         |
|        home          | The home site or repository site |   N     |   Y         |
|    source.mode       | The mode of the plugin. `local_so` for local `so` file; `remote_git` for remote source repository | Y  |   Y         |
|    source.path       | The `so` file path or the remote git repositry | Y |   Y |
| http_services.driver | The name of the http service driver which is used to enable the http services   | Y | N |
| http_services.routes | A route list to map the service endpoints to the plugin method with labels | Y | N |
| http_services.routes[i].route | The service endpoint definition        | Y | N |
| http_services.routes[i].method| The http method to apply on the route | Y | N |
| http_services.routes[i].label | Add the label to the plugin context when calling the plugin entry method | Y | N |

## Plugin Management

The following sample code shows how to manage the plugins with go-plugin.

```go
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
```

## Develop a plugin

* Step 1: Implement `func Execute(ctx context.PluginContext) error`
* Step 2: Build the plugin with command `go build -buildmode=plugin -o sample.so sample.go`

**NOTES:**
* If the logic is loop logic in goroutine, please make sure there is an exit way by listening the context done() channel or your goroutine may escape away when unloading the plugin
* If you need to handle multiple scenarios, just call your sub logic based on some values which are passed by plugin context. e.g:

```go
func Execute(ctx context.PluginContext) error {
    label := ctx.GetValue("label")
    switch label {
        case "case1":
          return call_subMethod1()
        case "case 2":
          return call_subMethod12()
        default:
          return errors.New("not suppotred")
    }
}
```

## Next steps

- [] Package the `plugin.json` and the `so` file as single `*.plg` file (with gzip)
- [] Build the plugin from source code @git repo
- [] Monitor and detect the plugin change in the plugin base dir
- [] Plugin hot upgrade
- [] Support http service onboarding drivers (beego first)
- [] Load plugins from internet
- [] Provide plugin metrics
- [] Enable API and run as rest services
- [] Provide GUI for plugin management
- [] Provided hooks for the lifecycle of plugin
- [] Add test cases and setup CI/CD

## Issues

* Currently plugins are only supported on Linux and macOS (Extened from golang plugin)
* Memory issues if discard a loaded `so`