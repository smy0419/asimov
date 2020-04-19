// Copyright (c) 2018-2020. The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package net

import (
	"errors"
	"fmt"
	"github.com/AsimovNetwork/asimov/socks"
	"net"
	"os"
	"strings"
	"time"
)

const (
	defaultConnectTimeout        = time.Second * 30
)

// NetAdapter is an interface representing a Net adapter traversal options for example Onion or
// Proxy. It provides methods to connect to server and lookup host.
type NetAdapter interface {
	// Connects to a server, with default connect timeout 30s
	// @params
	//   addr: server address in net type
	// @return
	//   Conn: a successful connection
	//   error: some detail information when connect failed.
	DialTimeout(addr net.Addr) (net.Conn, error)

	// Lookup resolves the IP of the given host using the correct DNS lookup
	// function depending on the configuration options.  For example, addresses will
	// be resolved using tor when the --proxy flag was specified unless --noonion
	// was also specified in which case the normal system DNS resolver will be used.
	//
	// Any attempt to resolve a tor address (.onion) will return an error since they
	// are not intended to be resolved outside of the tor proxy.
    Lookup(host string) ([]net.IP, error)

	// whether it support onion
	SupportOnion() bool
}

type fnetAdapter struct {
	proxy *socks.Proxy
	onionProxy *socks.Proxy
	noOnion bool
	timeout time.Duration
}

// Connects to a server
// @params
//   network: type of net
//   addr: server address
//   timeout: wait time of connect, if timeout, this function returns an error
// @return
//   Conn: a successful connection
//   error: some detail information when connect failed.
func (na * fnetAdapter) DialTimeout(addr net.Addr) (net.Conn, error) {
	if strings.Contains(addr.String(), ".onion:") {
		// noonion means the onion address dial function results in
		// an error.
		if na.noOnion {
			return nil, errors.New("tor has been disabled")
		}
		if na.onionProxy != nil {
			return na.onionProxy.DialTimeout(addr.Network(), addr.String(), na.timeout)
		}
	}

	if na.proxy != nil {
		return na.proxy.DialTimeout(addr.Network(), addr.String(), na.timeout)
	}

	return net.DialTimeout(addr.Network(), addr.String(), na.timeout)
}

// Lookup resolves the IP of the given host using the correct DNS lookup
// function depending on the configuration options.  For example, addresses will
// be resolved using tor when the --proxy flag was specified unless --noonion
// was also specified in which case the normal system DNS resolver will be used.
//
// Any attempt to resolve a tor address (.onion) will return an error since they
// are not intended to be resolved outside of the tor proxy.
func (na * fnetAdapter) Lookup(host string) ([]net.IP, error) {
	if strings.HasSuffix(host, ".onion") {
		return nil, fmt.Errorf("attempt to resolve tor address %s", host)
	}
	if na.proxy != nil {
		if na.onionProxy != nil {
			return torLookupIP(host, na.onionProxy.Addr)
		}
		if !na.noOnion {
			return torLookupIP(host, na.proxy.Addr)
		}
	}
	return net.LookupIP(host)
}

func (na * fnetAdapter) SupportOnion() bool {
	return !na.noOnion
}

// NewNetAdapter create a new fnetAdapter instance.
// Setup dial and DNS resolution (lookup) functions depending on the
// specified options.  The default is to use the standard
// net.DialTimeout function as well as the system DNS resolver.  When a
// proxy is specified, the dial function is set to the proxy specific
// dial function and the lookup is set to use tor (unless --noonion is
// specified in which case the system DNS resolver is used).
func NewNetAdapter(proxy, user, pass string, onionProxy, onionUser, onionPass string,
	torIsolation bool, noOnion bool) NetAdapter {
	nap := &fnetAdapter {
		noOnion: noOnion,
		timeout: defaultConnectTimeout,
	}
	if proxy != "" {
		// Tor isolation flag means proxy credentials will be overridden
		// unless there is also an onion proxy configured in which case
		// that one will be overridden.
		tmptorIsolation := false
		if torIsolation && onionProxy == "" &&
			(user != "" || pass != "") {
			tmptorIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified proxy user credentials")
		}
		nap.proxy = &socks.Proxy{
			Addr:         proxy,
			Username:     user,
			Password:     pass,
			TorIsolation: tmptorIsolation,
		}
	}
	if onionProxy != "" {
		nap.onionProxy = &socks.Proxy{
			Addr:         onionProxy,
			Username:     onionUser,
			Password:     onionPass,
			TorIsolation: torIsolation,
		}
	}

	return nap
}
