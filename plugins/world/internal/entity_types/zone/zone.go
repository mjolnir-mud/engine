package zone

type zoneType struct{}

func (z *zoneType) Name() string {
	return "zone"
}

func (z *zoneType) Create(_ string, args map[string]interface{}) map[string]interface{} {
	return args
}

var Type = &zoneType{}
