// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package protos

import (
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/common/serialization"
	"io"
	"net"
	"time"
)

// Services 8 bytes + ip 16 bytes + port 2 bytes + timestamp 4 bytes.
const maxNetAddressPayload = 30

// NetAddress defines information about a peer on the network including the time
// it was last seen, the services it supports, its IP address, and port.
type NetAddress struct {
	// Last time the address was seen.  This is, unfortunately, encoded as a
	// int64 on the protos and therefore is limited to 2106.  This field is
	// present in the flow version message (MsgVersion).
	Timestamp time.Time

	// Bitfield which identifies the services supported by the address.
	Services common.ServiceFlag

	// IP address of the peer.
	IP net.IP

	// Port the peer is using.  This is encoded in big endian on the protos
	// which differs from most everything else.
	Port uint16
}

// HasService returns whether the specified service is supported by the address.
func (na *NetAddress) HasService(service common.ServiceFlag) bool {
	return na.Services&service == service
}

// AddService adds service as a supported service by the peer generating the
// message.
func (na *NetAddress) AddService(service common.ServiceFlag) {
	na.Services |= service
}

// NewNetAddressIPPort returns a new NetAddress using the provided IP, port, and
// supported services with defaults for the remaining fields.
func NewNetAddressIPPort(ip net.IP, port uint16, services common.ServiceFlag) *NetAddress {
	return NewNetAddressTimestamp(time.Now(), services, ip, port)
}

// NewNetAddressTimestamp returns a new NetAddress using the provided
// timestamp, IP, port, and supported services. The timestamp is rounded to
// single second precision.
func NewNetAddressTimestamp(
	timestamp time.Time, services common.ServiceFlag, ip net.IP, port uint16) *NetAddress {
	// Limit the timestamp to one second precision since the protocol
	// doesn't support better.
	na := NetAddress{
		Timestamp: timestamp,
		Services:  services,
		IP:        ip,
		Port:      port,
	}
	return &na
}

// NewNetAddress returns a new NetAddress using the provided TCP address and
// supported services with defaults for the remaining fields.
func NewNetAddress(addr *net.TCPAddr, services common.ServiceFlag) *NetAddress {
	return NewNetAddressIPPort(addr.IP, uint16(addr.Port), services)
}

// readNetAddress reads an encoded NetAddress from r depending on the protocol
// version and whether or not the timestamp is included per ts.  Some messages
// like version do not include the timestamp.
func readNetAddress(r io.Reader, pver uint32, na *NetAddress) error {
	uint64time := uint64(0)
	if err := serialization.ReadUint64(r, &uint64time); err != nil {
		return err
	}
	na.Timestamp = time.Unix(int64(uint64time), 0)
	if err := serialization.ReadUint64(r, (*uint64)(&na.Services)); err != nil {
		return err
	}
	var ip [16]byte

	if err := serialization.ReadNBytes(r, ip[:], 16); err != nil {
		return err
	}

	var port uint16
	err := serialization.ReadUint16(r, &port)
	if err != nil {
		return err
	}

	*na = NetAddress{
		Timestamp: na.Timestamp,
		Services:  na.Services,
		IP:        ip[:],
		Port:      port,
	}
	return nil
}

// writeNetAddress serializes a NetAddress to w depending on the protocol
// version and whether or not the timestamp is included per ts.  Some messages
// like version do not include the timestamp.
func writeNetAddress(w io.Writer, pver uint32, na *NetAddress) error {
	// Ensure to always write 16 bytes even if the ip is nil.
	var ip [16]byte
	if na.IP != nil {
		copy(ip[:], na.IP.To16())
	}
	if err := serialization.WriteUint64(w, uint64(na.Timestamp.Unix())); err != nil {
		return err
	}
	if err := serialization.WriteUint64(w, uint64(na.Services)); err != nil {
		return err
	}
	if err := serialization.WriteNBytes(w, ip[:]); err != nil {
		return err
	}

	return serialization.WriteUint16(w, na.Port)
}
