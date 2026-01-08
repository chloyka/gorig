package effects

import (
	"reflect"
)

type InterpretedEffect struct {
	name      string
	enabled   bool
	processFn func([]float32)
}

func newInterpretedEffect(name string, enabled bool, processFn reflect.Value) *InterpretedEffect {
	return &InterpretedEffect{
		name:    name,
		enabled: enabled,
		processFn: func(samples []float32) {
			processFn.Call([]reflect.Value{reflect.ValueOf(samples)})
		},
	}
}

func (e *InterpretedEffect) Name() string {
	return e.name
}

func (e *InterpretedEffect) IsEnabled() bool {
	return e.enabled
}

func (e *InterpretedEffect) SetEnabled(enabled bool) {
	e.enabled = enabled
}

func (e *InterpretedEffect) Process(samples []float32) {
	e.processFn(samples)
}
