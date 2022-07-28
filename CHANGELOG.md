## Next Release
### Engine

#### New Features

* [template & theming engine plugin]()
  * new templates created for existing commands

#### Breaking Changes
* redis connectivity has been moved out of the world plugin into the engine
* plugin `Init` has been renamed to `Start`
* `engine.AddCommand` has been renamed to `engine.AddCLICommand`
* consolidated plugins into monorepo

### World

#### New Features
* move command can optionally not broadcast a message to members in a room

#### Bug Fixes

* miss naming of `room.MoveWithMesssageForSession()` to `room.MoveWithMessageForSession()`

