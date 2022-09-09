package instance

// BeforeStart - called before the engine starts.
var beforeStartCallbacks []func()

// BeforeStartForEnv - called before the engine starts for a specific environment
var beforeStartCallbacksForEnv map[string][]func()

// BeforeServiceStart - called before a service starts
var beforeServiceStartCallbacks map[string][]func()

// BeforeServiceStartForEnv - called before a service starts for a specific environment
var beforeServiceStartCallbacksForEnv map[string]map[string][]func()

func initializeBeforeStartCallbacks() {
	beforeStartCallbacks = make([]func(), 0)
	beforeStartCallbacksForEnv = make(map[string][]func())
	beforeServiceStartCallbacks = make(map[string][]func())
	beforeServiceStartCallbacksForEnv = make(map[string]map[string][]func())
}

func RegisterBeforeStartCallback(callback func()) {
	beforeStartCallbacks = append(beforeStartCallbacks, callback)
}

func RegisterBeforeStartCallbackForEnv(env string, callback func()) {
	_, ok := beforeStartCallbacksForEnv[env]

	if !ok {
		beforeStartCallbacksForEnv[env] = make([]func(), 0)
	}

	beforeStartCallbacksForEnv[env] = append(beforeStartCallbacksForEnv[env], callback)
}

func RegisterBeforeServiceStartCallback(service string, callback func()) {
	_, ok := beforeServiceStartCallbacks[service]

	if !ok {
		beforeServiceStartCallbacks[service] = make([]func(), 0)
	}

	beforeServiceStartCallbacks[service] = append(beforeServiceStartCallbacks[service], callback)
}

func RegisterBeforeServiceStartCallbackForEnv(service string, env string, callback func()) {
	_, ok := beforeServiceStartCallbacksForEnv[service]

	if !ok {
		beforeServiceStartCallbacksForEnv[service] = make(map[string][]func())
	}

	_, ok = beforeServiceStartCallbacksForEnv[service][env]

	if !ok {
		beforeServiceStartCallbacksForEnv[service][env] = make([]func(), 0)
	}

	beforeServiceStartCallbacksForEnv[service][env] = append(beforeServiceStartCallbacksForEnv[service][env], callback)
}

func callBeforeStartCallbacks() {
	log.Info().Msg("Calling before start callbacks")
	for _, callback := range beforeStartCallbacks {
		callback()
	}
}

func callBeforeStartCallbacksForEnv(env string) {
	log.Info().Msg("Calling before start callbacks for env " + env)
	for _, callback := range beforeStartCallbacksForEnv[env] {
		callback()
	}
}

func callBeforeServiceStartCallbacks(service string) {
	log.Info().Msg("Calling before service start callbacks for " + service)
	for _, callback := range beforeServiceStartCallbacks[service] {
		callback()
	}
}

func callBeforeServiceStartCallbacksForEnv(service string, env string) {
	log.Info().Msg("Calling before service start callbacks for " + service + " and env " + env)
	for _, callback := range beforeServiceStartCallbacksForEnv[service][env] {
		callback()
	}
}
