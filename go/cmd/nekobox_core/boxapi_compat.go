package main

import (
	"context"
	"errors"
	"net"
	"net/http"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing/common/bufio"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

// Compatibility layer for boxapi functions to work with standard sing-box

var ErrNoOutbound = errors.New("no outbound available")

// DialContext dials a connection through the sing-box instance
func DialContext(ctx context.Context, instance *box.Box, network, address string) (net.Conn, error) {
	outbound := instance.Outbound().Default()
	if outbound == nil {
		return nil, ErrNoOutbound
	}
	
	destination, err := M.ParseSocksaddr(address)
	if err != nil {
		return nil, err
	}
	
	return outbound.NewConnection(ctx, net.Conn(nil), adapter.InboundContext{
		Destination: destination,
	})
}

// DialUDP dials a UDP connection through the sing-box instance  
func DialUDP(ctx context.Context, instance *box.Box) (net.PacketConn, error) {
	outbound := instance.Outbound().Default()
	if outbound == nil {
		return nil, ErrNoOutbound
	}
	
	conn, err := outbound.NewPacketConnection(ctx, nil, adapter.InboundContext{})
	if err != nil {
		return nil, err
	}
	
	return bufio.NewNATPacketConn(bufio.NewPacketConn(conn), M.Socksaddr{}, M.Socksaddr{}), nil
}

// CreateProxyHttpClient creates an HTTP client that routes through the sing-box instance
func CreateProxyHttpClient(instance *box.Box) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return DialContext(ctx, instance, network, addr)
			},
		},
	}
}
