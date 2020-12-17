package gilfoyle

import (
	"github.com/dreamvo/gilfoyle/logging"
)

var (
	Logger logging.ILogger
)

func init() {
	l, err := logging.NewLogger(Config.Settings.Debug, true)
	if err != nil {
		panic(err)
	}
	Logger = l
}
