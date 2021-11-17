package fakexip

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("fakexip", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'fakexip'
	if c.NextArg() {
		return plugin.Error("fakexip", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Fakexip{}
	})

	return nil
}
