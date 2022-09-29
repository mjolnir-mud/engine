package server

import (
	"bytes"
	"fmt"
	events2 "github.com/mjolnir-mud/engine/plugins/sessions/events"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/pkg/logger"
	"github.com/rs/zerolog"
)

type connection struct {
	uuid                  string
	conn                  net.Conn
	stop                  chan bool
	logger                zerolog.Logger
	connectedAt           int64
	lastInputAt           int64
	remoteAddr            string
	sessManagerStartedSub engine.Subscription
	connectionReaderSub   engine.Subscription
}

func newConnection(conn net.Conn) *connection {
	uid := uuid.New().String()

	return &connection{
		uuid: uid,
		conn: conn,
		stop: make(chan bool),
		logger: logger.Instance.
			With().
			Str("component", "connection").
			Str("uuid", uid).
			Logger(),
		connectedAt: time.Now().UTC().Unix(),
		lastInputAt: time.Now().UTC().Unix(),
		remoteAddr:  conn.RemoteAddr().String(),
	}
}

// Start starts the connection
func (c *connection) Start() {
	c.logger.Debug().Msg("starting connection")
	_, err := c.conn.Write([]byte("Mjolnir MUD Engine\n"))

	if err != nil {
		c.logger.Error().Err(err).Msg("error occurred writing")
		c.Stop()
	}

	go func() {
		c.logger.Trace().Msg("starting connection loop")
		for {
			buf := make([]byte, 1024)

			_, err := c.conn.Read(buf)

			if err != nil {
				c.logger.Error().Err(err).Msg("error reading from connection")
				c.Stop()

				break
			}

			trimmed := bytes.Trim(buf, "\x00")
			lines := strings.Split(string(trimmed), "\r\n")

			// send each line to the world
			for _, line := range lines {
				if len(line) > 0 {
					err = engine.Publish(events2.PlayerInputEvent{
						Id:   c.uuid,
						Line: line,
					})

					if err != nil {
						c.logger.Error().Err(err).Msg("error publishing PlayerInput event")
						c.Stop()
						break
					}
				}
			}
		}
	}()

	sub := engine.Subscribe(events2.SessionRegistryStartedEvent{}, func(_ event.EventPayload) {
		c.logger.Debug().Msg("handling session manager started event")
		c.assertSession()

	})

	defer sub.Stop()

	sub = engine.Subscribe(events2.PlayerOutputEvent{
		Id: c.uuid,
	}, func(e event.EventPayload) {
		ev := &events2.PlayerOutputEvent{}
		err := e.Unmarshal(ev)

		c.logger.Debug().Msg("handling write to connection event")
		_, err = c.conn.Write([]byte(ev.Line))

		if err != nil {
			c.logger.Error().Err(err).Msg("error writing to connection")
			c.Stop()
		}
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("error subscribing to event")
		c.Stop()
		return
	}

	c.connectionReaderSub = sub

	c.assertSession()

	<-c.stop

	fmt.Println("Connection closed")
}

// Stop stops the connection
func (c *connection) Stop() {
	c.logger.Debug().Msg("stopping connection")
	_ = c.conn.Close()

	_ = engine.Publish(events2.PlayerDisconnectedEvent{
		Id: c.uuid,
	})

	c.stop <- true
}

func (c *connection) assertSession() {
	c.logger.Trace().Msg("asserting session")

	err := engine.Publish(events2.PlayerConnectedEvent{
		Id: c.uuid,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("error publishing AssertSession event")
		c.Stop()
		return
	}

}
