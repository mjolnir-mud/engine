package location

//
//import "github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
//
//type location struct{}
//
//func (e location) Name() string {
//	return "expiration"
//}
//
//func (e location) Component() string {
//	return "location"
//}
//
//func (e location) Match(key string, _ interface{}) bool {
//	return true
//}
//
//func (e location) WorldStarted() {}
//
//func (e location) ComponentAdded(entityId string, _ string, _ interface{}) error { return nil }
//
//func (e location) ComponentUpdated(entityId string, _ string, _ interface{}, _ interface{}) error {
//	return nil
//}
//
//func (e location) ComponentRemoved(entityId string, _ string, _ interface{}) error { return nil }
//
//func (e location) MatchingComponentAdded(entityId string, _ string, _ interface{}) error { return nil }
//
//func (e location) MatchingComponentUpdated(entityId string, _ string, _ interface{}, _ interface{}) error {
//	return nil
//}
//
//func (e location) MatchingComponentRemoved(entityId string, _ string, _ interface{}) error {
//	return nil
//}
//
//func AtLocation(entityId string) []string {
//	return entity_registry.AllEntitiesByTypeWithComponentValue("character", "location", entityId)
//}
//
//func AtLocationByType(entityId string, entityType string) []string {
//	return entity_registry.AllEntitiesByTypeWithComponentValue(entityType, "location", entityId)
//}
//
//func Set(entityId string, locationId string) {
//	entity_registry.SetStringComponent(entityId, "location", locationId)
//}
//
//var System = location{}
