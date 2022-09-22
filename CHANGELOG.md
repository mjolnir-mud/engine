## Next Release

### Engine

#### Bug Fixes
* Fix bug that prevents tests from correctly running due to a missing environment.

### DataSources

#### New Features
* Data Sources now can create a new entity by calling the `CreateEntity` or the `CreateEntityWithId` functions. It will 
* automatically set the entity metadata and return the entity.

#### Breaking Changes
* Data Sources interface has been updated to include an `AppendMetadata` function. This function will be called when
  an entity is created. It will be passed the entity metadata and is expected to return any metadata required by the 
  data source appended.


* The `Save` function has been changed to `SaveWithId`. A new `Save` function has been added that will automatically
  generate a new entity ID and call `SaveWithId` with the new ID, returning the new id.

### MongoDataSource

### Bug Fixes
* Fix bug that incorrectly handled ids preventing find functions from working correctly.

## 0.2.2
* Initial public release.