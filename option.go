package hvue

import (
	"github.com/gopherjs/gopherwasm/js"
	// "github.com/gopherjs/gopherjs/js"
)

// Config is the config object for NewVM.
type Config struct {
	js.Value
	// Data       js.Value `js:"data"`
	// Props      js.Value `js:"props"`
	// El         string   `js:"el"`
	// Methods    js.Value `js:"methods"`
	// Template   string   `js:"template"`
	// Computed   js.Value `js:"computed"`
	// Components js.Value `js:"components"`
	// Filters    js.Value `js:"filters"`

	// Not sure how to handle this yet
	// dataValue reflect.Value

	// Setters js.Value `js:"hvue_setters"`
}

func (c *Config) Data() js.Value       { return c.Get("data") }
func (c *Config) Props() js.Value      { return c.Get("props") }
func (c *Config) El() string           { return c.Get("el").String() }
func (c *Config) Methods() js.Value    { return c.Get("methods") }
func (c *Config) Template() string     { return c.Get("template").String() }
func (c *Config) Computed() js.Value   { return c.Get("computed") }
func (c *Config) Components() js.Value { return c.Get("components") }
func (c *Config) Filters() js.Value    { return c.Get("filters") }
func (c *Config) Setters() js.Value    { return c.Get("hvue_setters") }

func (c *Config) SetData(new js.Value)       { c.Set("data", new) }
func (c *Config) SetProps(new js.Value)      { c.Set("props", new) }
func (c *Config) SetEl(new string)           { c.Set("el", new) }
func (c *Config) SetMethods(new js.Value)    { c.Set("methods", new) }
func (c *Config) SetTemplate(new string)     { c.Set("template", new) }
func (c *Config) SetComputed(new js.Value)   { c.Set("computed", new) }
func (c *Config) SetComponents(new js.Value) { c.Set("components", new) }
func (c *Config) SetFilters(new js.Value)    { c.Set("filters", new) }
func (c *Config) SetSetters(new js.Value)    { c.Set("hvue_setters", new) }

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
	typ       js.Value    `js:"type"`
	required  bool        `js:"required"`
	def       interface{} `js:"default"`
	validator js.Value    `js:"validator"`
}

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
