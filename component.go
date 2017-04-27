package hvue

import "github.com/gopherjs/gopherjs/js"

func NewComponent(name string, opts ...option) {
	c := &Config{Object: NewObject()}
	c.Option(opts...)
	js.Global.Get("Vue").Call("component", name, c.Object)
}

func Props(props ...string) option {
	return func(c *Config) {
		if c.Props == js.Undefined {
			c.Props = NewArray()
		}
		for i, prop := range props {
			c.Props.SetIndex(i, prop)
		}
	}
}

func PropsO(props ...string) option {
	return func(c *Config) {
		if c.Props != js.Undefined {
			panic("Cannot use Props and PropsO in the same component")
		}
		// Do the rest ...
	}
}

func Template(template string) option {
	return func(c *Config) {
		c.Template = template
	}
}
