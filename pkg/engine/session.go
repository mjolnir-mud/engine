package engine

import "github.com/rs/zerolog"

type session struct {
	uuid        string
	connectedAt int64
	lastInputAt int64
	remoteAddr  string
	logger      zerolog.Logger
}
