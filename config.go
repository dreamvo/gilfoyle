package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/config"
)

var (
	Config config.Config
)

// NewConfig creates a new config object
// and load values from environment variables or config file.
// File paths can be both relative and absolute.
func NewConfig(files ...string) (*config.Config, error) {
	c, err := config.NewConfig(files...)
	if err != nil {
		return nil, err
	}

	Config = *c

	return &Config, err
}
