package main

import (
	"cnsguy/proxychainsd/log"
	"cnsguy/proxychainsd/socks"
	"errors"
	"fmt"
	"io"
	"net"
)

type Client struct {
	net.Conn
	*Server
	*log.PrefixedLogger
}

type ErrDisabledSocksVersion uint8

func (e ErrDisabledSocksVersion) Error() string {
	return fmt.Sprintf("Socks version %d support is disabled in config", e)
}

func NewClient(conn net.Conn, server *Server) *Client {
	return &Client{
		conn,
		server,
		log.NewLogger("[client] %s", conn.RemoteAddr().String(),
	)}
}

func tunnelLoop(log *log.PrefixedLogger, from net.Conn, to net.Conn) {
	buf := make([]byte, 65536)

	for {
		RefreshTimeout(from, 120) // XXX read from config
		n, err := from.Read(buf)

		if errors.Is(err, net.ErrClosed) {
			return
		} else if errors.Is(err, io.EOF) {
			log.Log("Disconnected")
			return
		} else if err != nil {
			log.Log("Misc error: %s", err.Error())
			return
		} else if n == 0 {
			log.Log("Disconnected")
			return
		}

		to.Write(buf[:n])
	}
}

func MakeTunnel() (net.Conn, error) {
	return net.Dial("tcp", "127.0.0.1:9050")
}

func (c *Client) Run() {
	defer c.Close()
	cmd, err := socks.ReadCommand(c)

	if err != nil {
		c.Log("Unacceptable socks command: %s", err.Error())
		return
	} else if cmd.Version() == 4 && !c.Server.ListenConf.EnableSocks4 {
		c.Log("Client tried to use socks4 but it's disabled in config")
		return
	} else if cmd.Version() == 5 && !c.Server.ListenConf.EnableSocks5 {
		c.Log("Client tried to use socks5 but it's disabled in config")
		return
	}

	c.Log("Got socks command %s (protocol version %d)", cmd.Type().String(), cmd.Version())
	t, err := MakeTunnel()

	if err != nil {
		c.Log("Failed to create tunnel: %s", err.Error())
		return
	}

	err = cmd.Send(t)

	if err != nil {
		c.Log("Failed to write command to tunnel socket: %s", err.Error())
		return
	}

	rsp, err := socks.ReadResponse(t, cmd.Version())

	if err != nil {
		c.Log("Error reading socks response: %s", err.Error())
		return
	}

	rsp.Send(c)

	go func() {
		defer t.Close()
		tunnelLoop(c.PrefixedLogger.Extend("[c->s]"), c, t)
	}()

	tunnelLoop(c.PrefixedLogger.Extend("[s->c]"), t, c)
}
