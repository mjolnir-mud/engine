package login

//
//import (
//	"github.com/mjolnir-mud/engine/plugins/world/internal/account"
//	"github.com/mjolnir-mud/engine/plugins/world/pkg/db"
//	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"golang.org/x/crypto/bcrypt"
//)
//
//// controller is the login controller, responsible handling user logins.
//type controller struct{}
//
//func (l controller) Name() string {
//	return "login"
//}
//
//func (l controller) Start(session session.Session) error {
//	_ = promptLoginUsername(session)
//
//	return nil
//}
//
//func (l controller) Resume(session session.Session) error {
//	return nil
//}
//
//func (l controller) Stop(session session.Session) error {
//	return nil
//}
//
//func (l controller) HandleInput(session session.Session, input string) error {
//	switch input {
//	case "new":
//		session.SetController("new_account")
//		return nil
//	default:
//		return handleInput(session, input)
//	}
//}
//
//func handleInput(session session.Session, input string) error {
//	switch session.GetIntFromFlash("step") {
//	case 1:
//		return handleUsername(session, input)
//	case 2:
//		return handlePassword(session, input)
//	}
//
//	return nil
//}
//
//func handleUsername(session session.Session, input string) error {
//	if input == "new" {
//		session.SetController("new_account")
//		return nil
//	}
//
//	session.SetInFlash("username", input)
//
//	return promptLoginPassword(session)
//}
//
//func handlePassword(session session.Session, input string) error {
//	username := session.GetStringFromFlash("username")
//
//	account := &account.Account{}
//	err := collection().FindOne(session.Context(), bson.M{"username": username}).Decode(account)
//
//	if err != nil {
//		session.WriteToConnection("Invalid username or password.")
//		return promptLoginUsername(session)
//	}
//
//	err = bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(input))
//
//	if err != nil {
//		session.WriteToConnection("Invalid username or password.")
//		return promptLoginUsername(session)
//	}
//
//	session.SetInStore("accountID", account.ID)
//	session.SetController("select_character")
//
//	return nil
//}
//
//func promptLoginPassword(session session.Session) error {
//	session.SetInFlash("step", 2)
//	session.WriteToConnection("Enter your password:")
//
//	return nil
//}
//
//func promptLoginUsername(session session.Session) error {
//	session.SetInFlash("step", 1)
//	session.WriteToConnection("Enter a username, or 'new' to create a new account:")
//
//	return nil
//}
//
//func collection() *mongo.Collection {
//	return db.Collection("accounts")
//}
//
//var Controller = controller{}
