package spec

import (
	"github.com/steven-zou/go-plugin/pkg/context"
)

/**
 * This is just a sample plugin to describe the interface of the plugin.
 * Each plugin should have an exported method named 'Execute' with a parameter
 * 'ctx: pkg.PluginContext'.
 *
 * Execute is the entry method of the plugin.
 */

//Execute the plugin logic here
func Execute(ctx context.PluginContext) error {
	return nil
}
