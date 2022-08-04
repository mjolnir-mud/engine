package system_registry

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	redis2 "github.com/go-redis/redis/v9"
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/constants"
	"github.com/mjolnir-mud/engine/plugins/ecs/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/ecs/pkg/system"
)

type registry struct {
	systems   map[string]system.System
	listeners map[string]subscription
}

type subscription struct {
	pubSub *redis2.PubSub
	stop   chan bool
}

func (s *subscription) Stop() {
	_ = s.pubSub.Close()
	s.stop <- true
}

// Start starts the registry.
func Start() {
	for _, s := range r.systems {
		startComponentListener(s)
	}
}

func Stop() {
	for _, s := range r.systems {
		stopComponentListener(s)
	}
}

// Register registers a system with the registry. If a system with the same name is already registered, it will be
// overwritten.
func Register(system system.System) {
	log.Info().Msgf("registering system %s", system.Name())
	r.systems[system.Name()] = system
}

func startComponentListener(s system.System) {
	log.Info().Msgf("starting component listener for system %s", s.Name())
	startComponentSetListener(s)
}

func startComponentSetListener(s system.System) {
	log.Info().Msgf("starting component set listener for system %s", s.Name())
	pubsub := engine.Redis.PSubscribe(context.Background(), keySpaceEventForSystem(s))

	sub := subscription{
		pubSub: pubsub,
		stop:   make(chan bool),
	}

	r.listeners[s.Name()] = sub

	go func() {
		for {
			select {
			case <-sub.stop:
				return
			case msg := <-pubsub.Channel():
				log.Trace().Msgf("received message %s", msg.Payload)
				parts := strings.Split(msg.Channel, ":")
				id := parts[1]
				name := parts[2]
				key := fmt.Sprintf("%s:%s", id, name)

				// if the id starts with __ don't call any callbacks
				if strings.HasPrefix(id, "__") {
					continue
				}

				switch msg.Payload {
				case "set":
					callSetCallbacks(s, id, key)
				case "hset":
					callHSetCallbacks(s, id, key)
				case "hdel":
					callHSetCallbacks(s, id, key)
				case "sadd":
					callSAddCallbacks(s, id, key)
				case "srem":
					callSAddCallbacks(s, id, key)
				case "del":
					callDelCallbacks(s, id, key)
				}
			}
		}
	}()

}

func callComponentAddedCallbacks(s system.System, id string, key string, value interface{}) {
	k := strings.Replace(key, fmt.Sprintf("%s:", id), "", 1)
	setComponentMeta(id, k, value)

	for _, sys := range r.systems {
		log.Trace().Msgf("calling component %s added callbacks for system %s", k, sys.Name())
		err := sys.ComponentAdded(id, k, value)

		if err != nil {
			log.Error().Err(err).Msgf("error calling component added callbacks for system %s", sys.Name())
		}
	}

	if s.Match(key, value) {
		log.Trace().Msgf("component %s added to system %s", k, s.Name())
		err := s.MatchingComponentAdded(id, k, value)

		if err != nil {
			log.Error().Err(err).Msgf("error calling matching component added callbacks for system %s", s.Name())
		}
	}
}

func callComponentUpdatedCallbacks(s system.System, id string, key string, oldValue interface{}, newValue interface{}) {
	k := strings.Replace(key, fmt.Sprintf("%s:", id), "", 1)
	for _, sys := range r.systems {
		log.Trace().Msgf("calling component %s updated callbacks for system %s", k, sys.Name())
		err := sys.ComponentUpdated(id, k, oldValue, newValue)

		if err != nil {
			log.Error().Err(err).Msgf("error calling component updated callbacks for system %s", sys.Name())
		}
	}

	if s.Match(key, newValue) {
		log.Trace().Msgf("component %s updated in system %s", k, s.Name())
		err := s.MatchingComponentUpdated(id, k, oldValue, newValue)

		if err != nil {
			log.Error().Err(err).Msgf("error calling matching component updated callbacks for system %s", s.Name())
		}
	}
	setComponentMeta(id, k, newValue)
}

func callDelCallbacks(s system.System, id string, key string) {
	var value interface{}
	k := strings.Replace(key, fmt.Sprintf("%s:", id), "", 1)

	for _, sys := range r.systems {
		log.Trace().Msgf("calling component deleted callbacks for system %s", sys.Name())
		valueType := engine.Redis.Get(context.Background(), fmt.Sprintf("%s:%s", constants.ComponentTypePrefix, key)).Val()

		switch valueType {
		case "string":
			value = engine.Redis.Get(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
		case "int":
			value = engine.Redis.Get(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
		case "int64":
			value = engine.Redis.Get(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
		case "map":
			m := engine.Redis.HGetAll(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
			value = make(map[string]interface{})

			for key, v := range m {
				value.(map[string]interface{})[key] = v
			}
		case "set":
			s := engine.Redis.SMembers(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
			value = make([]interface{}, len(s))

			for i, v := range s {
				value.([]interface{})[i] = v
			}
		}

		err := sys.ComponentRemoved(id, k)

		if err != nil {
			log.Error().Err(err).Msgf("error calling component deleted callbacks for system %s", sys.Name())
		}
	}

	if s.Match(k, nil) {
		log.Trace().Msgf("component %s deleted from system %s", k, s.Name())
		err := s.MatchingComponentRemoved(id, k)

		if err != nil {
			log.Error().Err(err).
				Msgf("error calling matching component deleted callbacks for system %s", s.Name())
		}
	}

	metaKeys := engine.Redis.Keys(context.Background(), fmt.Sprintf("__*:%s", k)).Val()

	for _, metaKey := range metaKeys {
		engine.Redis.Del(context.Background(), metaKey)
	}
}

func callHSetCallbacks(s system.System, id string, key string) {
	exists := engine.Redis.Exists(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
	currentStringValue := engine.Redis.HGetAll(context.Background(), key).Val()

	currentValue := make(map[string]interface{})

	for k, v := range currentStringValue {
		currentValue[k] = v
	}

	if exists == 0 {
		callComponentAddedCallbacks(s, id, key, currentValue)
		return
	}

	prevValue := engine.Redis.HGetAll(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()

	prevValueMap := make(map[string]interface{})

	for k, v := range prevValue {
		prevValueMap[k] = v
	}

	if reflect.DeepEqual(prevValue, currentValue) {
		return
	}

	callComponentUpdatedCallbacks(s, id, key, prevValueMap, currentValue)
}

func callSAddCallbacks(s system.System, id string, key string) {
	exists := engine.Redis.Exists(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
	currentStringValue := engine.Redis.SMembers(context.Background(), key).Val()

	currentValue := make([]interface{}, len(currentStringValue))

	for i, v := range currentStringValue {
		currentValue[i] = v
	}

	if exists == 0 {
		callComponentAddedCallbacks(s, id, key, currentValue)
		return
	}

	prevValue := engine.Redis.SMembers(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()

	prevValueSlice := make([]interface{}, len(prevValue))

	for i, v := range prevValue {
		prevValueSlice[i] = v
	}

	if reflect.DeepEqual(prevValueSlice, currentValue) {
		return
	}

	callComponentUpdatedCallbacks(s, id, key, prevValueSlice, currentValue)
}

func callSetCallbacks(s system.System, id string, key string) {
	exists := engine.Redis.Exists(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()
	currentValue := engine.Redis.Get(context.Background(), key).Val()

	if exists == 0 {
		callComponentAddedCallbacks(s, id, key, currentValue)
		return
	}

	prevValue := engine.Redis.Get(context.Background(), fmt.Sprintf("%s:%s", constants.PreviousValuePrefix, key)).Val()

	if prevValue == currentValue {
		return
	}

	callComponentUpdatedCallbacks(s, id, key, prevValue, currentValue)
}

func keySpaceEventForSystem(s system.System) string {
	return fmt.Sprintf("__keyspace@%d__:*:%s", 0, s.Component())
}

func setComponentMeta(id string, key string, value interface{}) {
	valueType := reflect.TypeOf(value).String()

	switch valueType {
	case "set":
		for _, i := range value.([]interface{}) {
			engine.Redis.SAdd(context.Background(), fmt.Sprintf("%s:%s:%s", constants.PreviousValuePrefix, id, key), i)
		}
	case "map":
		for k, v := range value.(map[string]interface{}) {
			engine.Redis.HSet(context.Background(), fmt.Sprintf("%s:%s:%s", constants.PreviousValuePrefix, id, key), k, v)
		}
	default:
		engine.Redis.Set(context.Background(), fmt.Sprintf("%s:%s:%s", constants.PreviousValuePrefix, id, key), value, 0)
	}

}

func stopComponentListener(s system.System) {
	log.Info().Msgf("stopping component listener for system %s", s.Name())
	stopComponentSetListener(s)
}

func stopComponentSetListener(s system.System) {
	log.Info().Msgf("stopping component set listener for system %s", s.Name())
	sub := r.listeners[s.Name()]
	sub.Stop()
	delete(r.listeners, s.Name())
}

var r = &registry{
	systems:   make(map[string]system.System),
	listeners: make(map[string]subscription),
}

var log = logger.Logger.
	With().
	Str("service", "systemRegistry").
	Logger()
