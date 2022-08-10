package server

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mjolnir-mud/engine"
	events2 "github.com/mjolnir-mud/engine/pkg/events"
	"github.com/mjolnir-mud/engine/plugins/telnet_portal/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/events"
	"github.com/nats-io/nats.go"
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
	sessManagerStartedSub *nats.Subscription
	connectionReaderSub   *nats.Subscription
}

func newConnection(conn net.Conn) *connection {
	uid := uuid.New().String()

	return &connection{
		uuid: uid,
		conn: conn,
		stop: make(chan bool),
		logger: logger.Logger.
			With().
			Str("service", "server").
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
					err = engine.PublishEvent(events2.SendToWorldTopic(c.uuid), &events2.SendToWorld{
						Line: line,
					})

					if err != nil {
						c.logger.Error().Err(err).Msg("error publishing SendToWorld event")
						c.Stop()
						break
					}
				}
			}

			if err != nil {
				fmt.Println("Error writing:", err.Error())
				c.Stop()

				break
			}
		}
	}()

	sub, err := engine.SubscribeToEvent(events.SessionManagerStartedTopic(), func(_ *events.SessionManagerStarted) {
		c.logger.Debug().Msg("handling session manager started event")
		c.assertSession()

	})

	if err != nil {
		c.logger.Error().Err(err).Msg("error subscribing to event")
		c.Stop()
		return
	}

	c.sessManagerStartedSub = sub

	sub, err = engine.SubscribeToEvent(events2.WriteToConnectionTopic(c.uuid), func(event *events2.WriteToConnection) {
		c.logger.Debug().Msg("handling write to connection event")
		_, err := c.conn.Write([]byte(event.Line))

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
	_ = c.sessManagerStartedSub.Unsubscribe()
	_ = c.connectionReaderSub.Unsubscribe()

	c.stop <- true
}

func (c *connection) assertSession() {
	c.logger.Trace().Msg("asserting session")

	err := engine.PublishEvent(events.AssertSessionTopic(), &events.AssertSession{
		UUID:        c.uuid,
		ConnectedAt: c.connectedAt,
		LastInputAt: c.lastInputAt,
		RemoteAddr:  c.remoteAddr,
	})

	if err != nil {
		c.logger.Error().Err(err).Msg("error publishing AssertSession event")
		c.Stop()
		return
	}

}
