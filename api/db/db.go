package db

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/ent"
)

var Client *ent.Client

// InitClient initialize a database connection
func InitClient(config config.DatabaseConfig) (err error) {
	datasource := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.Host,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	)

	Client, err = ent.Open(config.Dialect, datasource)

	return err
}
