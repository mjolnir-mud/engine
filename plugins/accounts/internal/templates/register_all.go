package templates

import templates2 "github.com/mjolnir-mud/engine/plugins/templates"

func RegisterAll() {
	templates2.RegisterTemplate(PromptUsernameTemplate)
	templates2.RegisterTemplate(PromptPasswordTemplate)
	templates2.RegisterTemplate(PromptEmailTemplate)
	templates2.RegisterTemplate(PromptNewPasswordTemplate)
	templates2.RegisterTemplate(PromptPasswordConfirmationTemplate)
	templates2.RegisterTemplate(PromptNewUsernameTemplate)
	templates2.RegisterTemplate(PromptNewEmailTemplate)
	templates2.RegisterTemplate(InvalidEmailAddressTemplate)
	templates2.RegisterTemplate(PasswordMatchFailTemplate)
	templates2.RegisterTemplate(UsernameTakenTemplate)
	templates2.RegisterTemplate(LoginInvalidTemplate)
}
