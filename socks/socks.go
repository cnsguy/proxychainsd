package socks

import (
	"bufio"
	"fmt"
	"io"
)

type ErrInvalidVersion uint8

type SocksCommand interface {
	Version() SocksVersion
	Type() SocksCommandType
	Send(io.Writer) error
}

type SocksResponse interface {
	Version() SocksVersion
	Send(io.Writer) error
	IsSuccess() bool
}

type SocksVersion uint8

const (
	Socks4 SocksVersion = 4
	Socks5              = iota
)

type SocksCommandType uint8

const (
	Connect SocksCommandType = iota
	Bind
)

func (e ErrInvalidVersion) Error() string {
	return fmt.Sprintf("Invalid socks version %d", e)
}

func (t SocksCommandType) String() string {
	switch t {
	case Connect:
		return "Connect"
	case Bind:
		return "Bind"
	default:
		return "Invalid SocksCommandType value"
	}
}

func ReadCommand(r io.Reader) (SocksCommand, error) {
	buf := bufio.NewReader(r)
	pkt, err := buf.Peek(1)

	if err != nil {
		return nil, err
	}

	ver := SocksVersion(pkt[0])

	switch ver {
	case Socks4:
		return ReadSocks4Command(buf)
	case Socks5:
		return ReadSocks5Command(buf)
	default:
		return nil, ErrInvalidVersion(ver)
	}
}

func ReadResponse(r io.Reader, ver SocksVersion) (SocksResponse, error) {
	switch ver {
	case Socks4:
		return ReadSocks4Response(r)
	case Socks5:
		return ReadSocks5Response(r)
	default:
		return nil, ErrInvalidVersion(ver)
	}
}
