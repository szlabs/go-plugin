package spec

import "github.com/steven-zou/go-plugin/pkg/context"

//PluginExecutor is the executor of the plugin
type PluginExecutor func(ctx context.PluginContext) error
