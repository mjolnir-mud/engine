package select_controller

//
//import (
//	"context"
//	"fmt"
//	"strconv"
//	"strings"
//
//	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
//	"github.com/mjolnir-mud/engine/plugins/world/pkg/db"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//)
//
//type selectCharacter struct{}
//
//func (c selectCharacter) Name() string {
//	return "select_character"
//}
//
//func (c selectCharacter) Start(session reactor.Session) error {
//	if countCharacters(session) == 0 {
//		session.SetController("create_character")
//		return nil
//	}
//
//	return promptSelectUser(session)
//}
//
//func (c selectCharacter) Resume(session reactor.Session) error {
//	return nil
//}
//
//func (c selectCharacter) Stop(session reactor.Session) error {
//	return nil
//}
//
//func (c selectCharacter) HandleInput(session reactor.Session, input string) error {
//	switch input {
//	case "create":
//		session.SetController("create_character")
//		return nil
//	default:
//		characters, err := lookupCharacters(session)
//
//		if err != nil {
//			return err
//		}
//
//		// convert the input to an integer
//		idx, err := strconv.Atoi(input)
//
//		if err != nil {
//			session.WriteToConnection("That is not a valid selection.")
//			return promptSelectUser(session)
//		}
//
//		// check that the index is within the range of the character list
//		if idx < 1 || idx > len(characters) {
//			session.WriteToConnection("That is not a valid selection.")
//			return promptSelectUser(session)
//		}
//
//		// set the character ID
//		id := characters[idx-1]["_id"].(string)
//		session.SetStringInStore("characterID", id)
//		entity_registry.SetStringComponent(id, "sessionID", session.ID())
//		entity_registry.SetStringComponent(id, "location", "limbo_prime")
//
//		session.SetController("game")
//
//		return nil
//
//	}
//}
//
//func promptSelectUser(session reactor.Session) error {
//	characters, err := lookupCharacters(session)
//
//	if err != nil {
//		return err
//	}
//
//	// map character components to map of strings
//	characterList := make([]map[string]string, 0)
//
//	for _, character := range characters {
//		components := make(map[string]string)
//		for key, value := range character["components"].(bson.M) {
//			components[key] = value.(string)
//		}
//
//		characterList = append(characterList, components)
//	}
//
//	// build user prompt list,
//	promptList := make([]string, 0)
//	promptList = append(promptList, "Select a character:")
//
//	for idx, character := range characterList {
//		promptList = append(promptList, fmt.Sprintf("%d. %s", idx+1, character["name"]))
//	}
//
//	// prompt user
//	session.WriteToConnection(strings.Join(promptList, "\n"))
//
//	return nil
//}
//
//func lookupCharacters(session reactor.Session) ([]bson.M, error) {
//	query := bson.M{"components.accountID": session.GetStringFromStore("accountID")}
//
//	cursor, err := collection().Find(context.Background(), query)
//
//	if err != nil {
//		return nil, err
//	}
//
//	var characters []bson.M
//
//	err = cursor.All(context.Background(), &characters)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return characters, nil
//}
//
//func countCharacters(session reactor.Session) int64 {
//	query := bson.M{"components.accountID": session.GetStringFromStore("accountID")}
//
//	count, err := collection().CountDocuments(context.Background(), query)
//
//	if err != nil {
//		return 0
//	}
//
//	return count
//}
//
//func collection() *mongo.Collection {
//	return db.Collection("characters")
//}
//
//var Name = selectCharacter{}
