package character

type characterType struct{}

func (c characterType) Name() string {
	return "character"
}

func (c characterType) Create(_ string, args map[string]interface{}) map[string]interface{} {
	args["persist"] = "characters"

	return args
}

var Type = characterType{}
