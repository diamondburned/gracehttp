package gracehttp

import (
	"fmt"
	"testing"
)

func TestParseAddr(t *testing.T) {
	type test struct {
		in   string
		net  string
		addr string
	}

	var tests = []test{
		{"/path/to.socket", "unix", "/path/to.socket"},
		{"unix://path/to.socket", "unix", "path/to.socket"},
		{"unix:///path/to.socket", "unix", "/path/to.socket"},
		{":20449", "tcp", ":20449"},
		{"0.0.0.0:20449", "tcp", "0.0.0.0:20449"},
		{"http://domainname.com", "tcp", "domainname.com"},
		{"tcp4://192.168.1.120:2048", "tcp4", "192.168.1.120:2048"},
		{
			"tcp6://24a3:c8ab:f219:e0d6:a631:a443:b001:5eee",
			"tcp6", "24a3:c8ab:f219:e0d6:a631:a443:b001:5eee",
		},
		{
			"tcp6://[24a3:c8ab:f219:e0d6:a631:a443:b001:5eee]:8081",
			"tcp6", "[24a3:c8ab:f219:e0d6:a631:a443:b001:5eee]:8081",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%s_%d", test.net, i), func(t *testing.T) {
			net, addr, err := parseAddr(test.in)
			if err != nil {
				t.Fatalf("failed to parse %q: %v", test.in, err)
			}

			if test.net != net {
				t.Errorf("net: expected %q, got %q", test.net, net)
			}
			if test.addr != addr {
				t.Errorf("addr: expected %q, got %q", test.addr, addr)
			}
		})
	}
}
