package db

import (
	"fmt"
	"github.com/dreamvo/gilfoyle/config"
	"github.com/dreamvo/gilfoyle/ent"
)

var Client *ent.Client

func NewClient(config *config.Config) (*ent.Client, error) {
	datasource := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		config.Services.DB.Host,
		config.Services.DB.Port,
		config.Services.DB.User,
		config.Services.DB.Database,
		config.Services.DB.Password,
	)

	Client, err := ent.Open(config.Services.DB.Dialect, datasource)

	return Client, err
}
