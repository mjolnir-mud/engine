package entities

type accountType struct{}

func (a accountType) Name() string {
	return "account"
}

func (a accountType) Create(args map[string]interface{}) map[string]interface{} {
	return args
}

var Account = &accountType{}
