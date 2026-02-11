package main

import (
	"context"
	"errors"
	"net"
	"net/http"

	box "github.com/sagernet/sing-box"
	M "github.com/sagernet/sing/common/metadata"
)

// Compatibility layer for boxapi functions to work with standard sing-box

var ErrNoOutbound = errors.New("no outbound available")

// DialContext dials a connection through the sing-box instance
func DialContext(ctx context.Context, instance *box.Box, network, address string) (net.Conn, error) {
	outbound := instance.Outbound().Default()
	if outbound == nil {
		return nil, ErrNoOutbound
	}
	
	destination := M.ParseSocksaddr(address)
	return outbound.DialContext(ctx, network, destination)
}

// DialUDP dials a UDP connection through the sing-box instance  
func DialUDP(ctx context.Context, instance *box.Box) (net.PacketConn, error) {
	outbound := instance.Outbound().Default()
	if outbound == nil {
		return nil, ErrNoOutbound
	}
	
	// Use ListenPacket to create a packet connection
	return outbound.ListenPacket(ctx, M.Socksaddr{})
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
