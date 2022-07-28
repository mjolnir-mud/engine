package templates

import "fmt"

type WalkingContext struct {
	Direction string
	Focus string
	Name string
}

type walkingTemplate struct {}

func (t walkingTemplate) Name() string {
	return "walking"
}

func (t walkingTemplate) Style() string {
	return "walks"
}

func (t walkingTemplate) Render(ctx interface{}) (string, error) {
	c := ctx.(*WalkingContext)

	switch c.Focus {
	case "self":
		return fmt.Sprintf("You walk %s.", c.Direction), nil
	case "other":
		return fmt.Sprintf("%s walks %s.", c.Name, c.Direction), nil
	case "no-exit":
		return fmt.Sprintf("There is no exit to the %s.", c.Direction), nil
	}

	return "", fmt.Errorf("invalid focus%s", c.Focus)
}

var Walking = walkingTemplate{}