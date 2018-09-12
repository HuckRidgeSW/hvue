package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherwasm/js"
)

// Config is the config object for NewVM.
type Config struct {
	js.Value

	DataType  js.Type
	dataValue reflect.Value
}

// Data and DataFunc both return the same underlying slot.
func (c *Config) Data() js.Value     { return c.Get("data") }
func (c *Config) DataFunc() js.Value { return c.Data() }

func (c *Config) Props() js.Value      { return c.Get("props") }
func (c *Config) El() string           { return c.Get("el").String() }
func (c *Config) Methods() js.Value    { return c.Get("methods") }
func (c *Config) Template() string     { return c.Get("template").String() }
func (c *Config) Computed() js.Value   { return c.Get("computed") }
func (c *Config) Components() js.Value { return c.Get("components") }
func (c *Config) Filters() js.Value    { return c.Get("filters") }
func (c *Config) Setters() js.Value    { return c.Get("hvue_setters") }
func (c *Config) Watchers() js.Value   { return c.Get("watch") }

// SetData and SetDataFunc both set the same underlying slot.
func (c *Config) SetData(new js.Value) {
	if new.Type() != js.TypeObject {
		panic("SetData must use an object; got " + new.Type().String() + ", value " + new.String())
	}
	c.Set("data", new)
	c.DataType = js.TypeObject
}

// The type of function passed to SetDataFunc, to initialize the fields for a
// new data object in a Vue component.
type DataFuncT func(vm *VM, dataObj js.Value) interface{}

func (c *Config) SetDataFunc(newF DataFuncT, fieldNames ...string) {
	templateObj := NewObject()
	for _, v := range fieldNames {
		templateObj.Set(v, "")
	}
	// data needs to be a real JS function that returns a real JS value.
	// wasm_new_data_func returns such a function; said function also calls the
	// newF callback to initialize the data slots at a later time.
	c.Set("data",
		js.Global().Call("wasm_new_data_func",
			templateObj,
			js.NewCallback(func(args []js.Value) {
				// Initialize the new data object; get a Go object back.
				value := newF(&VM{Value: args[0]}, args[1])
				storeDataID(args[1], value, c)

				// FIXME: This should be able to return a value
				// return value
			})))
	c.DataType = js.TypeFunction
}

func (c *Config) SetProps(new js.Value)      { c.Set("props", new) }
func (c *Config) SetEl(new string)           { c.Set("el", new) }
func (c *Config) SetMethods(new js.Value)    { c.Set("methods", new) }
func (c *Config) SetTemplate(new string)     { c.Set("template", new) }
func (c *Config) SetComputed(new js.Value)   { c.Set("computed", new) }
func (c *Config) SetComponents(new js.Value) { c.Set("components", new) }
func (c *Config) SetFilters(new js.Value)    { c.Set("filters", new) }
func (c *Config) SetSetters(new js.Value)    { c.Set("hvue_setters", new) }
func (c *Config) SetWatchers(new js.Value)   { c.Set("watch", new) }

type ComponentOption func(*Config)

// Option sets the options specified.
func (c *Config) Option(opts ...ComponentOption) {
	for _, opt := range opts {
		opt(c)
	}
}

type PropOption func(*PropConfig)

// PropConfig is the config object for Props
type PropConfig struct {
	js.Value
	// typ       js.Value    `js:"type"`
	// required  bool        `js:"required"`
	// def       interface{} `js:"default"`
	// validator js.Value    `js:"validator"`
}

func (pc *PropConfig) SetType(t js.Value)       { pc.Set("type", t) }
func (pc *PropConfig) SetRequired(r bool)       { pc.Set("required", r) }
func (pc *PropConfig) SetDefault(d interface{}) { pc.Set("default", d) }

func (p *PropConfig) Option(opts ...PropOption) {
	for _, opt := range opts {
		opt(p)
	}
}

type pOptionType int

const (
	PString   pOptionType = iota
	PNumber               = iota
	PBoolean              = iota
	PFunction             = iota
	PObject               = iota
	PArray                = iota
	// Not sure how to do custom types yet
)

type DirectiveOption func(*DirectiveConfig)

// DirectiveConfig is the config object for configuring a directive.
type DirectiveConfig struct {
	js.Value
	// Bind             js.Value `js:"bind"`
	// Inserted         js.Value `js:"inserted"`
	// Update           js.Value `js:"update"`
	// ComponentUpdated js.Value `js:"componentUpdated"`
	// Unbind           js.Value `js:"unbind"`
	// Short            js.Value `js:"short"`
}

func (dc *DirectiveConfig) Short() js.Value       { return dc.Get("short") }
func (dc *DirectiveConfig) SetShort(new js.Value) { dc.Set("short", new) }

func (c *DirectiveConfig) Option(opts ...DirectiveOption) {
	for _, opt := range opts {
		opt(c)
	}
}
