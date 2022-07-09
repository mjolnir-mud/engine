package engine

import "github.com/spf13/viper"

func setConfig() {
	viper.SetEnvPrefix("mjolnir")
	err := viper.BindEnv("nats_url")

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
}
