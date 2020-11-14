package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/config"
	"github.com/jinzhu/configor"
)

var (
	Config config.Config
)

// NewConfig creates a new config object
// and load values from environment variables or config file.
// File paths can be both relative and absolute.
func NewConfig(files ...string) (*config.Config, error) {
	err := configor.Load(&Config, files...)

	return &Config, err
}
