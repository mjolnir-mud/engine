package server

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_Start(t *testing.T) {
	server := New()

	go server.Start()

	c, err := net.Dial("tcp", "localhost:2323")

	assert.Nilf(t, err, "Error dialing: %v", err)

	ch := make(chan []byte)

	go func() {
		buf := make([]byte, 1024)
		_, err = c.Read(buf)

		ch <- bytes.Trim(buf, "\x00")

		assert.Nilf(t, err, "Error reading: #{err}")
	}()

	select {
	case res := <-ch:
		assert.Equal(t, "Mjolnir MUD Engine\n", string(res), "invalid response")
	case <-time.After(5 * time.Second):
		t.Fail()
	}
}
