//go:generate go run github.com/swaggo/swag/cmd/swag init -g ./api.go --parseDependency
package api

import (
	_ "github.com/dreamvo/gilfoyle/api/docs"
	_ "github.com/swaggo/swag/gen"
	_ "github.com/urfave/cli"
	_ "github.com/urfave/cli/v2"
)
