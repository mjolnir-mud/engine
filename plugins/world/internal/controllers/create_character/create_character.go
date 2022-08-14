package create_character

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type selectCharacter struct{}

func (c selectCharacter) Name() string {
	return "create_character"
}

func (c selectCharacter) Start(session reactor.Session) error {
	session.WriteToConnection("What is your character's name?")

	return nil
}

func (c selectCharacter) Resume(_ reactor.Session) error {
	return nil
}

func (c selectCharacter) Stop(_ reactor.Session) error {
	return nil
}

func (c selectCharacter) HandleInput(session reactor.Session, input string) error {
	validation := validateCharacterName(input)

	if validation != nil {
		session.WriteToConnection(validation.Error())
		return c.Start(session)
	}

	id := strings.ToLower(input)
	count, err := collection().CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	if count > 0 {
		session.WriteToConnection("That character name is already taken.")
		return c.Start(session)
	}

	accountID := session.GetStringFromStore("accountID")

	err = entity_registry.Add("character", id, map[string]interface{}{
		"name":      input,
		"accountID": accountID,
	})

	if err != nil {
		return err
	}

	session.SetStringInStore("characterID", id)

	session.SetController("game")
	return nil
}

// validateCharacterName validates the name of the character. It must be at least 4 characters long and  must be made
// up of only alphabetical characters.
func validateCharacterName(name string) error {
	if len(name) < 4 {
		//goland:noinspection GoErrorStringFormat
		return errors.New("Character name must be at least 4 characters long.")
	}

	if !regexp.MustCompile("^[a-zA-Z]+$").MatchString(name) {
		//goland:noinspection GoErrorStringFormat
		return errors.New("Character name must be made up of only alphabetical characters.")
	}

	return nil
}

func collection() *mongo.Collection {
	return db.Collection("characters")
}

var Controller = &selectCharacter{}
