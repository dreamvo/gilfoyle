package db

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/ent"
)

var Client *ent.Client

// InitClient init the public Client variable
func InitClient(config *config.Config) error {
	datasource := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.Services.DB.Host,
		config.Services.DB.Port,
		config.Services.DB.User,
		config.Services.DB.Database,
		config.Services.DB.Password,
	)

	c, err := ent.Open(config.Services.DB.Dialect, datasource)

	Client = c

	return err
}
