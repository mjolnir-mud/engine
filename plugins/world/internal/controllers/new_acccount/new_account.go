package new_acccount

import (
	"net/mail"

	"github.com/mjolnir-mud/engine/plugins/world/internal/account"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"

	"golang.org/x/crypto/bcrypt"

	"github.com/mjolnir-mud/engine/plugins/world/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type controller struct{}

var Controller = controller{}

func (n controller) Name() string {
	return "new_account"
}

func (n controller) Start(session session.Session) error {
	return promptSigninUsername(session)
}

func (n controller) Resume(session session.Session) error {
	return nil
}

func (n controller) Stop(session session.Session) error {
	return nil
}

func (n controller) HandleInput(session session.Session, input string) error {
	switch session.GetIntFromFlash("step") {
	case 1:
		return handleUsername(session, input)
	case 2:
		return handleEmail(session, input)
	case 3:
		return handlePassword(session, input)
	case 4:
		return handlePasswordConfirmation(session, input)
	}

	return nil
}

func handlePassword(session session.Session, input string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input), 4)

	if err != nil {
		return err
	}

	session.SetInFlash("password", string(hash))

	return promptPasswordConfirmation(session)
}

func handleEmail(session session.Session, input string) error {
	_, err := mail.ParseAddress(input)

	if err != nil {
		session.WriteToConnection("Invalid email address.")
		return promptEmail(session)
	}

	session.SetInFlash("email", input)
	return promptPassword(session)
}

func handleUsername(session session.Session, input string) error {
	c, err := collection().CountDocuments(session.Context(), bson.D{{"username", input}})

	if err != nil {
		return err
	}

	if c > 0 {
		session.WriteToConnection("That username is already taken.")
		return promptSigninUsername(session)
	}

	session.SetInFlash("username", input)
	return promptEmail(session)
}

func handlePasswordConfirmation(session session.Session, input string) error {
	hashedPassword := session.GetStringFromFlash("password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input))

	if err != nil {
		session.WriteToConnection("Passwords do not match.")
		return promptPassword(session)
	}

	_, err = collection().InsertOne(session.Context(), account.Account{
		Username:       session.GetStringFromFlash("username"),
		Email:          session.GetStringFromFlash("email"),
		HashedPassword: hashedPassword,
	})

	if err != nil {
		return err
	}

	session.WriteToConnection("Account created.")

	return nil
}

func promptSigninUsername(session session.Session) error {
	session.SetInFlash("step", 1)
	session.WriteToConnection("Enter a username:")

	return nil
}

func promptEmail(session session.Session) error {
	session.SetInFlash("step", 2)
	session.WriteToConnection("Enter an email address:")

	return nil
}

func promptPassword(session session.Session) error {
	session.SetInFlash("step", 3)
	session.WriteToConnection("Enter a password:")

	return nil
}

func promptPasswordConfirmation(session session.Session) error {
	session.SetInFlash("step", 4)
	session.WriteToConnection("Confirm your password:")

	return nil
}

func collection() *mongo.Collection {
	return db.Collection("accounts")
}
