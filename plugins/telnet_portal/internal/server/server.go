package server

import (
	"fmt"
	"net"

	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/pkg/config"
	"github.com/rs/zerolog"
)

var log zerolog.Logger
var stop = make(chan bool)
var connections = map[string]*connection{}

func Start(cfg *config.Configuration) {
	log = logger.Instance.With().
		Str("component", "server").
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Logger()

	log.Info().Msg("starting server")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	if err != nil {
		panic(err)
	}

	log.Trace().Msg("starting server loop")

	go func() {
		for {
			log.Trace().Msg("waiting for connection")
			conn, err := listener.Accept()

			if err != nil {
				panic(err)
			}

			select {
			case <-stop:
				_ = listener.Close()
				break
			default:
				handleConnection(conn)
			}
		}
	}()
}

func Stop() {
	log.Info().Msg("stopping server")

	for _, c := range connections {
		c.Stop()
	}

	stop <- true
}

func handleConnection(conn net.Conn) {
	log.Trace().Msg("handling new connection")
	c := newConnection(conn)
	connections[c.uuid] = c

	go c.Start()
}
