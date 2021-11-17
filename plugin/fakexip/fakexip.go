// Package fakexip implements a plugin like xip, return a reverse ip addr
package fakexip

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

const name = "fakexip"

// Fakexip is a plugin that parse the domain as xip and reverse it
// to CoreDNS.
type Fakexip struct{}

// ServeDNS implements the plugin.Handler interface.
func (wh Fakexip) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true
	var rr dns.RR
	ipv4RE := regexp.MustCompile(`(^|_)(((25[0-5]|(2[0-4]|1?[0-9])?[0-9])_){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9]))($|_)`)
	switch state.Family() {
	case 1:
		if !ipv4RE.Match([]byte(state.QName())){
			fmt.Println("pattern mismatch")
			return dns.RcodeFormatError, fmt.Errorf("format error")
		}
		match := ipv4RE.FindStringSubmatch(state.QName())[2]
		parts := make([]string, 4)
		for i, v := range strings.Split(match, "_"){
			parts[3-i] = v
		}
		ip := strings.Join(parts, ".")
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.(*dns.A).A = net.ParseIP(ip).To4()
	default:
		return dns.RcodeBadMode, fmt.Errorf("only A type supported")
	}

	a.Answer = []dns.RR{rr}

	w.WriteMsg(a)

	return 0, nil
}

// Name implements the Handler interface.
func (wh Fakexip) Name() string { return name }
