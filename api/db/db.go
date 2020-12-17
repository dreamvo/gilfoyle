package db

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/ent"
)

// NewClient returns a new database connection
func NewClient(config config.DatabaseConfig) (*ent.Client, error) {
	datasource := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.Host,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	)

	return ent.Open(config.Dialect, datasource)
}
