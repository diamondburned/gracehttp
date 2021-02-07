package gracehttp

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ListenAddr parses the given address and returns either a TCP or Unix
// listener. Below lists acceptable formats.
//
//    unixpacket://relative/path/to.socket
//    unix:///path/to.socket
//    /path/to.socket
//
//    http://address
//    tcp4://address
//    tcp6://address
//    tcp://address
//    address
//
func ListenAddr(addr string) (net.Listener, error) {
	return ListenAddrCfg(context.Background(), addr, net.ListenConfig{})
}

// ListenAddrCfg is ListenAddr with additional parameters.
func ListenAddrCfg(ctx context.Context, addr string, cfg net.ListenConfig) (net.Listener, error) {
	network, address, err := parseAddr(addr)
	if err != nil {
		return nil, err
	}

	// Ensure that the socket is cleaned up because we're not gracefully
	// handling closes.
	if strings.HasPrefix(network, "unix") {
		if err := os.Remove(address); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, errors.Wrap(err, "failed to clean up old socket")
			}
		}
	}

	var lcfg net.ListenConfig

	l, err := lcfg.Listen(ctx, network, address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen")
	}

	return l, nil
}

func parseAddr(addr string) (network, address string, err error) {
	parts := strings.SplitN(addr, "://", 2)

	if len(parts) == 2 {
		address = parts[1]
		network = parts[0]
	} else {
		address = parts[0]

		if strings.HasPrefix(address, "/") {
			network = "unix"
		} else {
			network = "tcp"
		}
	}

	switch network {
	case "tcp4", "tcp6", "tcp", "unix", "unixpacket":
		// acceptable scheme
	case "http":
		network = "tcp"
	default:
		return "", "", fmt.Errorf("unknown scheme %s", network)
	}

	return
}
