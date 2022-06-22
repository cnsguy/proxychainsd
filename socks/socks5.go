package socks

import "io"

type XXX struct {
}

func (XXX) Error() string {
	return ""
}

func ReadSocks5Command(r io.Reader) (SocksCommand, error) {
	return nil, XXX{}
}

func ReadSocks5Response(r io.Reader) (SocksResponse, error) {
	return nil, XXX{}
}
