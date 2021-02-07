package gracehttp_test

import (
	"io"
	"net/http"
	"time"

	"github.com/diamondburned/gracehttp"
)

const addr = "unix:///tmp/gracehttp-example.sock"

func Example() {
	// ListenAndServeAsync will only block until the listener is being served.
	server := gracehttp.ListenAndServeAsync(addr, http.HandlerFunc(handle))
	defer server.ShutdownTimeout(5 * time.Second)

	time.Sleep(time.Second)
}

func handle(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, 世界")
}
