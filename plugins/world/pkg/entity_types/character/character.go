package character

import "github.com/spf13/viper"

type characterType struct{}

func (c characterType) Name() string {
	return "character"
}

func (c characterType) Create(_ string, args map[string]interface{}) map[string]interface{} {
	args["persist"] = "characters"

	start := viper.GetString("character_starting_location")

	if args["location"] == nil {
		args["location"] = start
	}

	return args
}

var Type = characterType{}
