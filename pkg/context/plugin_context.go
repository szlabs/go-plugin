package context

import (
	"context"
	"strings"
	"time"
)

//PluginContext help to provide related information/parameters to the
//plugin execution entry method.
//PluginContext inherits all from the context.Context
type PluginContext interface {
	context.Context
	ValueContext
}

//ValueContext defines behaviors of handling values with string keys
type ValueContext interface {
	//Get value by the key
	GetValue(key string) interface{}

	//Set value to context
	SetValue(key string, value interface{})
}

//BasePluginContext implemented as default plugin context
type BasePluginContext struct {
	//For compatible with system context
	basedOnContext context.Context

	//For keeping values
	valueMap map[string]interface{}
}

//GetValue implements 'GetValue' in ValueContext interface
func (bpc *BasePluginContext) GetValue(key string) interface{} {
	return bpc.valueMap[key]
}

//SetValue implements 'SetValue' in ValueContext interface
func (bpc *BasePluginContext) SetValue(key string, value interface{}) {
	if len(strings.TrimSpace(key)) > 0 {
		//nil value is allowed
		bpc.valueMap[key] = value
	}
}

//Deadline implements 'Deadline' in context.Context
func (bpc *BasePluginContext) Deadline() (deadline time.Time, ok bool) {
	return bpc.basedOnContext.Deadline()
}

//Done implements 'Done' in context.Context
func (bpc *BasePluginContext) Done() <-chan struct{} {
	return bpc.basedOnContext.Done()
}

//Err implements 'Err' in context.Context
func (bpc *BasePluginContext) Err() error {
	return bpc.basedOnContext.Err()
}

//Value implements 'Value' in context.Context
func (bpc *BasePluginContext) Value(key interface{}) interface{} {
	return bpc.basedOnContext.Value(key)
}

//Background build the base plugin context based on the context.
func Background() PluginContext {
	return &BasePluginContext{
		basedOnContext: context.Background(),
		valueMap:       make(map[string]interface{}),
	}
}
