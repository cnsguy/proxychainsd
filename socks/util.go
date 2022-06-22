package socks

import (
	"encoding/binary"
	"fmt"
	"net"
)

type ErrIPv6 struct{}
type ErrUnknownIPSize int

func (e ErrIPv6) Error() string {
	return "Can't encode an IPv6 address in an uint32"
}

func (e ErrUnknownIPSize) Error() string {
	return fmt.Sprintf("Unknown IP size %d", e)
}

func UInt32ToIP(ip uint32) net.IP {
	r := make(net.IP, 4)
	binary.BigEndian.PutUint32(r, ip)
	return r
}

func IPToUInt32(ip net.IP) (uint32, error) {
	l := len(ip)

	switch l {
	case 4:
		return binary.BigEndian.Uint32(ip), nil
	case 6:
		return 0, ErrIPv6{}
	default:
		return 0, ErrUnknownIPSize(l)
	}
}
