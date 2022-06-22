package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	respSuccess = 90
)

const (
	socks4Connect = 1
	socks4Bind    = 2
)

type Socks4CommandPacket struct {
	Version uint8
	Command uint8
	Port    uint16
	IP      uint32
	UserID  uint8
}

type Socks4ResponsePacket struct {
	Version uint8
	Result  uint8
	Port    uint16
	IP      uint32
}

type ErrSocks4UserIDSpecified struct{}
type ErrSocks4InvalidCommand uint8

func (*ErrSocks4UserIDSpecified) Error() string {
	return "Socks4: Specifying UserID is unsupported"
}

func (e ErrSocks4InvalidCommand) Error() string {
	return fmt.Sprintf("Socks4: Invalid command %d", e)
}

type Socks4ConnectCommand struct {
	DestIP net.IP
	Port   uint16
}

func (c *Socks4ConnectCommand) Send(w io.Writer) error {
	ip, err := IPToUInt32(c.DestIP)

	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, Socks4CommandPacket{
		4,
		1,
		c.Port,
		ip,
		0,
	})
}

func (*Socks4ConnectCommand) Version() SocksVersion {
	return Socks4
}

func (*Socks4ConnectCommand) Type() SocksCommandType {
	return SocksCommandType(Connect)
}

type Socks4BindCommand struct {
	BindIP net.IP
	Port   uint16
}

func (c *Socks4BindCommand) Send(w io.Writer) error {
	panic("NYI")
}

func (*Socks4BindCommand) Version() SocksVersion {
	return Socks4
}

func (*Socks4BindCommand) Type() SocksCommandType {
	return SocksCommandType(Bind)
}

func ReadSocks4Command(r io.Reader) (SocksCommand, error) {
	var pkt Socks4CommandPacket
	err := binary.Read(r, binary.BigEndian, &pkt)

	if err != nil {
		return nil, err
	} else if pkt.UserID != 0 {
		return nil, &ErrSocks4UserIDSpecified{}
	}

	ip := UInt32ToIP(pkt.IP)

	switch pkt.Command {
	case socks4Connect:
		return &Socks4ConnectCommand{ip, pkt.Port}, nil
	case socks4Bind:
		return &Socks4BindCommand{ip, pkt.Port}, nil
	default:
		return nil, ErrSocks4InvalidCommand(pkt.Command)
	}
}

type Socks4Response struct {
	Result uint8
}

func (r *Socks4Response) IsSuccess() bool {
	return r.Result == 90
}

func (r *Socks4Response) Version() SocksVersion {
	return Socks4
}

func (s *Socks4Response) Send(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, Socks4ResponsePacket{
		0,
		s.Result,
		0,
		0,
	})
}

func ReadSocks4Response(r io.Reader) (SocksResponse, error) {
	var pkt Socks4ResponsePacket
	err := binary.Read(r, binary.BigEndian, &pkt)

	if err != nil {
		return nil, err
	}

	return &Socks4Response{pkt.Result}, nil
}
