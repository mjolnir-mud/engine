package area

type areaType struct{}

func (a *areaType) Name() string {
	return "area"
}

func (a *areaType) Create(_ string, args map[string]interface{}) map[string]interface{} {
	return args
}

var Type = &areaType{}
