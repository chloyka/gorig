package effects

type Effect interface {
	Process(samples []float32)
	Name() string
	IsEnabled() bool
}
