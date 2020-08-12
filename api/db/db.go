package db

import (
	"github.com/dreamvo/gilfoyle/ent"
)

var Client *ent.Client

func NewClient() (*ent.Client, error) {
	return ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
}
