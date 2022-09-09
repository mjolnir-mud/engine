package instance

// BeforeStop - called before the engine stops.
var beforeStopCallbacks []func()

// BeforeStopForEnv - called before the engine stops for a specific environment
var beforeStopCallbacksForEnv map[string][]func()

// BeforeServiceStop - called before a service stops
var beforeServiceStopCallbacks map[string][]func()

// BeforeServiceStopForEnv - called before a service stops for a specific environment
var beforeServiceStopCallbacksForEnv map[string]map[string][]func()

func initializeBeforeStopCallbacks() {
	beforeStopCallbacks = make([]func(), 0)
	beforeStopCallbacksForEnv = make(map[string][]func())
	beforeServiceStopCallbacks = make(map[string][]func())
	beforeServiceStopCallbacksForEnv = make(map[string]map[string][]func())
}

func RegisterBeforeStopCallback(callback func()) {
	beforeStopCallbacks = append(beforeStopCallbacks, callback)
}

func RegisterBeforeStopCallbackForEnv(env string, callback func()) {
	_, ok := beforeStopCallbacksForEnv[env]

	if !ok {
		beforeStopCallbacksForEnv[env] = make([]func(), 0)
	}

	beforeStopCallbacksForEnv[env] = append(beforeStopCallbacksForEnv[env], callback)
}

func RegisterBeforeServiceStopCallback(service string, callback func()) {
	_, ok := beforeServiceStopCallbacks[service]

	if !ok {
		beforeServiceStopCallbacks[service] = make([]func(), 0)
	}

	beforeServiceStopCallbacks[service] = append(beforeServiceStopCallbacks[service], callback)
}

func RegisterBeforeServiceStopCallbackForEnv(service string, env string, callback func()) {
	_, ok := beforeServiceStopCallbacksForEnv[service]

	if !ok {
		beforeServiceStopCallbacksForEnv[service] = make(map[string][]func())
	}

	_, ok = beforeServiceStopCallbacksForEnv[service][env]

	if !ok {
		beforeServiceStopCallbacksForEnv[service][env] = make([]func(), 0)
	}

	beforeServiceStopCallbacksForEnv[service][env] = append(beforeServiceStopCallbacksForEnv[service][env], callback)
}

func callBeforeStopCallbacks() {
	log.Info().Msg("Calling before stop callbacks")
	for _, callback := range beforeStopCallbacks {
		callback()
	}
}

func callBeforeStopCallbacksForEnv(env string) {
	log.Info().Msg("Calling before stop callbacks for env " + env)
	for _, callback := range beforeStopCallbacksForEnv[env] {
		callback()
	}
}

func callBeforeServiceStopCallbacks(service string) {
	log.Info().Msg("Calling before service stop callbacks for " + service)
	for _, callback := range beforeServiceStopCallbacks[service] {
		callback()
	}
}

func callBeforeServiceStopCallbacksForEnv(service string, env string) {
	log.Info().Msg("Calling before service stop callbacks for " + service + " and env " + env)
	for _, callback := range beforeServiceStopCallbacksForEnv[service][env] {
		callback()
	}
}
