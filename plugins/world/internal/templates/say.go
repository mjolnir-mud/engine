package templates

import "fmt"

type SayContext struct {
	Focus string
	Name string
	Message string
}

type template struct {}

func (t template) Name() string {
	return "say"
}

func (t template) Style() string {
	return "says"
}

func (t template) Render(ctx interface{}) (string, error) {
	c := ctx.(*SayContext)

	switch c.Focus {
	case "self":
		return fmt.Sprintf("You say, \"%s\"", c.Message), nil
	case "other":
		return fmt.Sprintf("%s says, \"%s\"", c.Name, c.Message), nil
	}

	return "", fmt.Errorf("invalid focus%s", c.Focus)
}

var Say = template{}