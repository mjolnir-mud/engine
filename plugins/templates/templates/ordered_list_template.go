package templates

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/templates/internal/theme_registry"
	"strings"
)

type orderedList struct{}

func (t orderedList) Name() string {
	return "ordered_list"
}

func (t orderedList) Style() string {
	return "default"
}

func (t orderedList) Render(list interface{}) (string, error) {
	listSlice := list.([]string)

	for i, item := range listSlice {
		itemNumber := i + 1
		str, err := theme_registry.Render("ordered-list-item", fmt.Sprintf("%d. ", itemNumber))

		if err != nil {
			return "", err
		}

		listSlice[i] = str + item
	}

	return fmt.Sprint(strings.Join(listSlice, "\n")), nil
}

var OrderedList = &orderedList{}
