/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package plugin

import (
	"context"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/mongo_data_source/internal/logger"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configuration struct {
	MongoURL string
	Database string
}

var Configs = map[string]func(c *Configuration) *Configuration{
	"default": func(c *Configuration) *Configuration {
		c.MongoURL = "mongodb://localhost:27017"
		c.Database = "mjolnir_dev"

		return c
	},
}

var Plugin = &plugin{}

var log zerolog.Logger

type plugin struct {
	database *mongo.Database
}

func (p *plugin) Name() string {
	return "mongo_data_source"
}

func (p *plugin) Registered() error {
	engine.RegisterBeforeServiceStartCallback("world", func() {
		logger.Start()
		log = logger.Instance
		env := engine.GetEnv()
		defaultConfig := Configs["default"](&Configuration{})
		config := Configs[env](defaultConfig)

		if config == nil {
			log.Fatal().Msg("no config for environment")
			panic("no config for environment")
		}

		log = logger.
			Instance.
			With().
			Str("mongo_url", config.MongoURL).
			Str("database", config.Database).
			Logger()

		log.Info().Msg("starting mongo connection")
		c, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURL))

		if err != nil {
			log.Fatal().Err(err).Msg("error connecting to mongo")
			panic(err)
		}

		err = c.Connect(context.Background())

		if err != nil {
			log.Fatal().Err(err).Msg("error connecting to mongo")
			panic(err)
		}

		p.database = c.Database(config.Database)

		log.Info().Msg("mongo connection started")
	})

	engine.RegisterBeforeStopCallback(func() {
		log.Info().Msg("stopping mongo connection")
		p.Disconnect()
	})

	return nil
}

func (p *plugin) Disconnect() {
	_ = p.database.Client().Disconnect(context.Background())
}

func (p *plugin) Drop() {
	_ = p.database.Drop(context.Background())
}

func (p *plugin) Collection(name string) *mongo.Collection {
	return p.database.Collection(name)
}
