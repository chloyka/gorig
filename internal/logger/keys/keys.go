package keys

import "go.uber.org/zap"

type StringKey func(value string) zap.Field

type IntKey func(value int) zap.Field

type BoolKey func(value bool) zap.Field

type Float64Key func(value float64) zap.Field

type StringsKey func(value []string) zap.Field

func String(key string) StringKey {
	return func(value string) zap.Field {
		return zap.String(key, value)
	}
}

func Int(key string) IntKey {
	return func(value int) zap.Field {
		return zap.Int(key, value)
	}
}

func Bool(key string) BoolKey {
	return func(value bool) zap.Field {
		return zap.Bool(key, value)
	}
}

func Float64(key string) Float64Key {
	return func(value float64) zap.Field {
		return zap.Float64(key, value)
	}
}

func Strings(key string) StringsKey {
	return func(value []string) zap.Field {
		return zap.Strings(key, value)
	}
}
