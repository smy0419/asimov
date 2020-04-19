// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package addrmgr

import (
	"net"
	"strings"
	"fmt"
	"runtime"
)

// OnionAddr implements the net.Addr interface and represents a tor address.
type OnionAddr struct {
	addr string
}

// String returns the onion address.
//
// This is part of the net.Addr interface.
func (oa *OnionAddr) String() string {
	return oa.addr
}

// Network returns "onion".
//
// This is part of the net.Addr interface.
func (oa *OnionAddr) Network() string {
	return "onion"
}

// Ensure OnionAddr implements the net.Addr interface.
var _ net.Addr = (*OnionAddr)(nil)

// SimpleAddr implements the net.Addr interface with two struct fields
type SimpleAddr struct {
	net, addr string
}

// String returns the address.
//
// This is part of the net.Addr interface.
func (a SimpleAddr) String() string {
	return a.addr
}

// Network returns the network.
//
// This is part of the net.Addr interface.
func (a SimpleAddr) Network() string {
	return a.net
}

// Ensure SimpleAddr implements the net.Addr interface.
var _ net.Addr = SimpleAddr{}

// parseListeners determines whether each listen address is IPv4 and IPv6 and
// returns a slice of appropriate net.Addrs to listen on with TCP. It also
// properly detects addresses which apply to "all interfaces" and adds the
// address as both IPv4 and IPv6.
func ParseListeners(addrs []string) ([]net.Addr, error) {
	netAddrs := make([]net.Addr, 0, len(addrs)*2)
	for _, addr := range addrs {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			// Shouldn't happen due to already being normalized.
			return nil, err
		}

		// Empty host or host of * on plan9 is both IPv4 and IPv6.
		if host == "" || (host == "*" && runtime.GOOS == "plan9") {
			netAddrs = append(netAddrs, SimpleAddr{net: "tcp4", addr: addr})
			netAddrs = append(netAddrs, SimpleAddr{net: "tcp6", addr: addr})
			continue
		}

		// Strip IPv6 zone id if present since net.ParseIP does not
		// handle it.
		zoneIndex := strings.LastIndex(host, "%")
		if zoneIndex > 0 {
			host = host[:zoneIndex]
		}

		// Parse the IP.
		ip := net.ParseIP(host)
		if ip == nil {
			return nil, fmt.Errorf("'%s' is not a valid IP address", host)
		}

		// To4 returns nil when the IP is not an IPv4 address, so use
		// this determine the address type.
		if ip.To4() == nil {
			netAddrs = append(netAddrs, SimpleAddr{net: "tcp6", addr: addr})
		} else {
			netAddrs = append(netAddrs, SimpleAddr{net: "tcp4", addr: addr})
		}
	}
	return netAddrs, nil
}