package main

import (
	"cnsguy/proxychainsd/log"
	"fmt"
	"net"
	"sync"
)

type ServerAddr struct { // XXX
	BindIP string
	Port   uint16
}

func (s *ServerAddr) String() string {
	return fmt.Sprintf("%s:%d", s.BindIP, s.Port)
}

type Server struct {
	ListenConf *ListenerConfig
	ChainConf []ChainConfig
	*log.PrefixedLogger
}

func NewServer(listenConf *ListenerConfig, chainConf []ChainConfig) *Server {
	return &Server{
		listenConf,
		chainConf,
		log.NewLogger("[server] %s:%d", listenConf.BindIP, listenConf.Port),
	}
}

func (s *Server) Run(wg sync.WaitGroup) {
	defer wg.Done()
	s.Log("Starting server (socks4 enabled: %t socks5 enabled: %t)", s.ListenConf.EnableSocks4, s.ListenConf.EnableSocks5)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.ListenConf.BindIP, s.ListenConf.Port))

	if err != nil {
		s.Log("Could not open listener: %s", err.Error())
		return
	}

	s.Log("Listener opened and awaiting connections")

	for {
		conn, err := listener.Accept()

		if err != nil {
			s.Log("Error during accept loop: %s", err.Error())
			break
		}

		s.Log("Accepted %s", conn.RemoteAddr().String())
		c := NewClient(conn, s)
		go c.Run()
	}
}
