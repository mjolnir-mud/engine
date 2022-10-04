package templates

import "github.com/mjolnir-mud/engine/plugins/templates/internal/template_registry"

func Register() {
	template_registry.Register(OrderedList)
}
