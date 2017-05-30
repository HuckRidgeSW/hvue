package hvue

import (
	"reflect"

	"github.com/gopherjs/gopherjs/js"
)

type Config struct {
	*js.Object
	El         string     `js:"el"`
	Data       *js.Object `js:"data"`
	Methods    *js.Object `js:"methods"`
	Props      *js.Object `js:"props"`
	Template   string     `js:"template"`
	Computed   *js.Object `js:"computed"`
	Components *js.Object `js:"components"`

	dataValue reflect.Value
}

type ComponentOption func(*Config)

// Option sets the options specified.
func (c *Config) Option(opts ...ComponentOption) {
	for _, opt := range opts {
		opt(c)
	}
}

type PropOption func(*PropConfig)

type PropConfig struct {
	*js.Object
	typ       *js.Object  `js:"type"`
	required  bool        `js:"required"`
	def       interface{} `js:"default"`
	validator *js.Object  `js:"validator"`
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

type DirectiveConfig struct {
	*js.Object
	Bind             *js.Object `js:"bind"`
	Inserted         *js.Object `js:"inserted"`
	Update           *js.Object `js:"update"`
	ComponentUpdated *js.Object `js:"componentUpdated"`
	Unbind           *js.Object `js:"unbind"`
	Short            *js.Object `js:"short"`
}

func (c *DirectiveConfig) Option(opts ...DirectiveOption) {
	for _, opt := range opts {
		opt(c)
	}
}
