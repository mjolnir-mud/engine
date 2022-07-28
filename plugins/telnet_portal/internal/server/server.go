package server

import (
	"fmt"
	"net"
	"os"

	"github.com/mjolnir-mud/engine/plugins/telnet_portal/internal/logger"
)

type server struct {
	listener        net.Listener
	startedHandlers []func()
	connections     map[string]*connection
}

var log = logger.Logger.
	With().
	Str("plugin", "telnet_portal").
	Str("service", "server").
	Logger()

// New creates a new Telnet server.
func New() *server {
	return &server{
		startedHandlers: make([]func(), 0),
		connections:     make(map[string]*connection),
	}
}

func (s *server) Start() {
	log.Info().Msg("starting server")
	listener, err := net.Listen("tcp", "localhost:2323")

	s.listener = listener

	if err != nil {
		fmt.Print(fmt.Errorf("error listening: %v", err))
		os.Exit(2)
	}

	defer s.Stop()

	log.Trace().Msg("calling started handlers")
	for _, handler := range s.startedHandlers {
		handler()
	}

	log.Trace().Msg("starting server loop")
	for {
		log.Trace().Msg("waiting for connection")
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		log.Trace().Msg("handling new connection")
		go s.handleConnection(conn)
	}
}

func (s *server) Stop() {
	log.Info().Msg("stopping server")
	for _, connection := range s.connections {
		log.Debug().Msgf("stopping connection: %s", connection.uuid)
		connection.Stop()
	}

	_ = s.listener.Close()
}

func (s *server) WhenStarted(f func()) {
	s.startedHandlers = append(s.startedHandlers, f)
}

func (s *server) handleConnection(conn net.Conn) {
	log.Trace().Msg("creating new connection")
	connection := newConnection(conn)
	s.connections[connection.uuid] = connection

	log.Debug().Msg("starting connection")
	go connection.Start()
}
