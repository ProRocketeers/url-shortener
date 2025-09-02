//go:build tools

package main

// this is essentially a module that only gets build with `tools` flag
// unused, just to stop Go from dropping the dependency when `go mod tidy`
// this is used in Atlas config file
import _ "ariga.io/atlas-provider-gorm/gormschema"
