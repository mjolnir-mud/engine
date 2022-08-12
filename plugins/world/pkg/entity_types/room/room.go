package room

type roomType struct{}

func (r *roomType) Name() string {
	return "room"
}

func (r *roomType) Create(_ string, args map[string]interface{}) map[string]interface{} {
	return args
}

var Type = &roomType{}
