package instance

// AfterStop - called after the engine stops.
var afterStopCallbacks []func()

// AfterStopForEnv - called after the engine stops for a specific environment
var afterStopCallbacksForEnv map[string][]func()

// AfterServiceStop - called after a service stops
var afterServiceStopCallbacks map[string][]func()

// AfterServiceStopForEnv - called after a service stops for a specific environment
var afterServiceStopCallbacksForEnv map[string]map[string][]func()

func initializeAfterStopCallbacks() {
	afterStopCallbacks = make([]func(), 0)
	afterStopCallbacksForEnv = make(map[string][]func())
	afterServiceStopCallbacks = make(map[string][]func())
	afterServiceStopCallbacksForEnv = make(map[string]map[string][]func())
}

func RegisterAfterStopCallback(callback func()) {
	afterStopCallbacks = append(afterStopCallbacks, callback)
}

func RegisterAfterStopCallbackForEnv(env string, callback func()) {
	afterStopCallbacksForEnv[env] = append(afterStopCallbacksForEnv[env], callback)
}

func RegisterAfterServiceStopCallback(service string, callback func()) {
	_, ok := afterServiceStopCallbacks[service]

	if !ok {
		afterServiceStopCallbacks[service] = make([]func(), 0)
	}

	afterServiceStopCallbacks[service] = append(afterServiceStopCallbacks[service], callback)
}

func RegisterAfterServiceStopCallbackForEnv(service string, env string, callback func()) {
	_, ok := afterServiceStopCallbacksForEnv[service]

	if !ok {
		afterServiceStopCallbacksForEnv[service] = make(map[string][]func())
	}

	_, ok = afterServiceStopCallbacksForEnv[service][env]

	if !ok {
		afterServiceStopCallbacksForEnv[service][env] = make([]func(), 0)
	}

	afterServiceStopCallbacksForEnv[service][env] = append(afterServiceStopCallbacksForEnv[service][env], callback)
}

func callAfterStopCallbacks() {
	log.Info().Msg("Calling after stop callbacks")
	for _, callback := range afterStopCallbacks {
		callback()
	}
}

func callAfterStopCallbacksForEnv(env string) {
	log.Info().Msg("Calling after stop callbacks for env " + env)
	for _, callback := range afterStopCallbacksForEnv[env] {
		callback()
	}
}

func callAfterServiceStopCallbacks(service string) {
	log.Info().Msg("Calling after service stop callbacks for " + service)
	for _, callback := range afterServiceStopCallbacks[service] {
		callback()
	}
}

func callAfterServiceStopCallbacksForEnv(service string, env string) {
	log.Info().Msg("Calling after service stop callbacks for " + service + " and env " + env)
	for _, callback := range afterServiceStopCallbacksForEnv[service][env] {
		callback()
	}
}
