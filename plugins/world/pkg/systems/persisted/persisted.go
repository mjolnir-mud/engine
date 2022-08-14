package persisted

//
//import (
//	"context"
//
//	"github.com/mjolnir-mud/engine/plugins/ecs"
//	"github.com/mjolnir-mud/engine/plugins/world/pkg/db"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//type system struct{}
//
//type entity struct {
//	ID         string                 `bson:"_id"`
//	Components map[string]interface{} `bson:"components"`
//}
//
//func (s system) Name() string {
//	return "persisted"
//}
//
//func (s system) Component() string {
//	return "persist"
//}
//
//func (s system) Match(key string, _ interface{}) bool {
//	return true
//}
//
//func (s system) WorldStarted() {
//}
//
//func (s system) ComponentAdded(entityId string, _ string, _ interface{}) error {
//	Persist(entityId)
//
//	return nil
//}
//
//func (s system) ComponentUpdated(entityId string, _ string, _ interface{}, _ interface{}) error {
//	Persist(entityId)
//
//	return nil
//}
//
//func (s system) ComponentRemoved(entityId string, _ string, _ interface{}) error {
//	Persist(entityId)
//
//	return nil
//}
//
//func (s system) MatchingComponentAdded(_ string, _ string, _ interface{}) error { return nil }
//
//func (s system) MatchingComponentUpdated(_ string, _ string, _ interface{}, _ interface{}) error {
//	return nil
//}
//
//func (s system) MatchingComponentRemoved(_ string, _ string, _ interface{}) error { return nil }
//
//func Persist(entityId string) {
//	collection, err := ecs.GetStringComponent(entityId, "persist")
//
//	if err != nil {
//		return
//	}
//
//	if collection != "" {
//		components := entity_registry.AllComponents(entityId)
//
//		e := entity{
//			ID:         entityId,
//			Components: components,
//		}
//
//		c := db.Collection(collection)
//
//		opts := &options.UpdateOptions{}
//		opts.SetUpsert(true)
//
//		_, err := c.UpdateOne(context.Background(), bson.M{"_id": entityId}, bson.M{"$set": e}, opts)
//
//		if err != nil {
//			panic(err)
//		}
//	}
//}
//
//var System = system{}
