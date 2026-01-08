package keys

import "go.uber.org/zap"

var Error = func(err error) zap.Field {
	return zap.Error(err)
}
