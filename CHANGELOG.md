## Next Release

### Accounts

#### New Features
* Password may not be either the username or email address

#### Breaking Changes
* Directories have been re-organized to be more consistent with the rest of the application

### Engine

#### Bug Fixes
* Fix bug that prevents tests from correctly running due to a missing environment.

### ECS

#### Breaking Changes
* The `EntityType` interface has been updated to require a `Validate` function. This function will becalled before an
  entity is added to the world. It should return an error if the entity is invalid, nil otherwise.


* The `EntityType` `Create` function has been renamed to `New` to better reflect its purpose.

* The `ecs.CreateEntity` function has been renamed to `ecs.NewEntity` to better reflect its purpose.

### DataSources

#### New Features
* Data Sources now can create a new entity by calling the `NewEntity` or the `NewEntityWithId` functions. It will 
  automatically set the entity metadata and return the entity.

#### Breaking Changes
* Data Sources interface has been updated to include an `AppendMetadata` function. This function will be called when
  an entity is created. It will be passed the entity metadata and is expected to return any metadata required by the 
  data source appended.

* The `Save` function has been changed to `SaveWithId`. A new `Save` function has been added that will automatically
  generate a new entity ID and call `SaveWithId` with the new ID, returning the new id.

* Restructured plugin package

### MongoDataSource

### Bug Fixes
* Fix bug that incorrectly handled ids preventing find functions from working correctly.

## 0.2.2
* Initial public release.