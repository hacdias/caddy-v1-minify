package minify

import (
	"strings"

	"github.com/mholt/caddy"
)

// Config is the configuration struct needed for this plugin
type config struct {
	Excludes []string
	BasePath string
}

func parse(c *caddy.Controller) (*config, error) {
	conf := &config{
		Excludes: []string{},
		BasePath: "/",
	}

	for c.Next() {
		args := c.RemainingArgs()

		switch len(args) {
		case 1:
			conf.BasePath = args[0]
			conf.BasePath = strings.TrimSuffix(conf.BasePath, "/")
			conf.BasePath += "/"
		}

		for c.NextBlock() {
			switch c.Val() {
			case "exclude":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				conf.Excludes = strings.Split(c.Val(), " ")
			}
		}
	}

	return conf, nil
}
