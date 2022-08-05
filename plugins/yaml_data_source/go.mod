module github.com/mjolnir-mud/engine/plugins/yaml_data_source

go 1.17

require (
	github.com/mjolnir-mud/engine/plugins/data_sources v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.27.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
)

replace github.com/mjolnir-mud/engine/plugins/data_sources => ../data_sources
