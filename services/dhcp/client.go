package dhcp

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func (ths *Client) SendPacket(data []byte) (*Packet, error) {
	_, err := ths.conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("can't send: %w", err)
	}
	buf := make([]byte, 2048)
	_, err = bufio.NewReader(ths.conn).Read(buf)
	if err != nil {
		return nil, fmt.Errorf("can't read response: %w", err)
	}
	pkt, err := NewPacket(buf)
	if err != nil {
		return nil, fmt.Errorf("can't parse response: %w", err)
	}
	return pkt, nil
}

func (ths *Client) Close() error {
	return ths.conn.Close()
}

func NewDhcpClient(address string) (*Client, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}
